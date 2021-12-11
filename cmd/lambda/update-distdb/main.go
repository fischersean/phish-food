package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	db "github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/etl"
	_ "github.com/fischersean/phish-food/internal/tzinit"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DB_NAME = "kettle.db"
)

func clearTbl(conn *sql.DB, tblName string) (err error) {
	clearStmtString := fmt.Sprintf("DELETE FROM %s", tblName)
	clearStmt, err := conn.Prepare(clearStmtString)
	if err != nil {
		return err
	}
	_, err = clearStmt.Exec()
	return err
}

func repopulateSymbols(conn *sql.DB, s3Service *s3.S3) (err error) {

	// Empty the symbols table
	err = clearTbl(conn, "Symbols")
	if err != nil {
		return err
	}

	tradeablesBucket := os.Getenv("TRADEABLES_BUCKET")
	tickerPopulation, err := etl.GetTradeableSecurities(s3Service, tradeablesBucket)
	if err != nil {
		return err
	}

	stmtString := "INSERT INTO Symbols (Ticker, Exchange, FullName, ETF) VALUES(?, ?, ?, ?)"
	stmt, err := conn.Prepare(stmtString)
	if err != nil {
		return err
	}

	for _, sym := range tickerPopulation {
		if _, err = stmt.Exec(sym.Symbol, sym.Exchange, sym.FullName, sym.ETF); err != nil {
			p := strings.Split(sym.Symbol, ":")
			if p[0] == "File Creation Time" || p[0] == "" {
				// reset
				err = nil
				continue
			}
			log.Println(sym.Symbol)
			return err
		}
	}

	return err

}

func repopulateCounts(conn *sql.DB, sess *session.Session) (err error) {

	dynamoConn, err := db.Connect(db.ConnectionInput{
		Session: sess,
	})
	if err != nil {
		return err
	}

	// get oldest date
	stmtString := `
		SELECT 
		FormatedDate
		FROM Counts
		ORDER BY FormatedDate DESC
		LIMIT 1`
	stmt, err := conn.Prepare(stmtString)
	if err != nil {
		return err
	}

	res, err := stmt.Query()
	if err != nil {
		return err
	}

	var strDate string // the formated date. will convert later
	res.Next()
	err = res.Scan(&strDate)
	if err != nil {
		return err
	}
	res.Close()

	stopDate, err := time.Parse("20060102", strDate)
	if err != nil {
		return err
	}

	// Starting from today and working backwards, fill in all missing dates
	for date := time.Now(); date.After(stopDate); date = date.Add(-24 * time.Hour) {
		for _, sub := range etl.FetchTargets {
			etlRecords, err := dynamoConn.GetEtlResultsRecord(db.EtlResultsQueryInput{
				Subreddit: sub,
				Date:      date,
			})
			if err != nil {
				return err
			}

			stmtString := `
				INSERT INTO Counts (
					Subreddit,
					FormatedDate,
					Hour,
					Ticker,
					PostScore,
					CommentScore,
					PostMentions,
					CommentMentions,
					TotalScore
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
			stmt, err := conn.Prepare(stmtString)
			if err != nil {
				return err
			}

			for _, record := range etlRecords {
				formatDate := strings.Split(record.Id, "_")
				for _, data := range record.Data {
					if _, err := stmt.Exec(sub,
						formatDate[1],
						record.Hour,
						data.Stock.Symbol,
						data.Count.PostScore,
						data.Count.CommentScore,
						data.Count.PostMentions,
						data.Count.CommentMentions,
						data.Count.TotalScore,
					); err != nil {
						log.Printf("%s: %s", err.Error(), record.Id)
						// Reset since we dont really care about fixing errors. We only want to know one happened
						err = nil
					}
				}
			}
		}
	}

	return err
}

func downloadDb(svc *s3.S3, distBucketName string, dbFname string) (err error) {

	// Download db file from bucket
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(distBucketName),
		Key:    aws.String(DB_NAME),
	},
	)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	var buf []byte
	buf, err = ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dbFname, buf, 0644)
	if err != nil {
		return err
	}

	return err
}

func uploadDb(svc *s3.S3, distBucketName string, dbFname string) (err error) {

	buf, err := ioutil.ReadFile(dbFname)
	if err != nil {
		return err
	}

	rs := bytes.NewReader(buf)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(distBucketName),
		Body:   rs,
		Key:    aws.String(DB_NAME),
	},
	)

	return err
}

func Handler() (err error) {

	// Init session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	S3Service := s3.New(sess)

	// TODO: I am handling this file poorly.
	// I need to fix this and adjust each func responsibilities
	var dbFname string
	if os.Getenv("DEV") == "YES" {
		dbFname = DB_NAME
	} else {
		dbFname = "/tmp/" + DB_NAME
	}

	distBucketName := os.Getenv("DIST_BUCKET")
	err = downloadDb(S3Service, distBucketName, dbFname)
	if err != nil {
		return err
	}

	if os.Getenv("DEV") != "YES" {
		defer os.Remove(dbFname)
	}

	// Perform db update
	conn, err := sql.Open("sqlite3", dbFname)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Updating symbols table")
	err = repopulateSymbols(conn, S3Service)
	if err != nil {
		return err
	}
	log.Println("Finished updating symbols table")

	log.Println("Updating counts table")
	err = repopulateCounts(conn, sess)
	if err != nil {
		return err
	}
	log.Println("Finished updating counts table")

	// Upload the file back to s3
	err = uploadDb(S3Service, distBucketName, dbFname)
	return err
}

func main() {
	if os.Getenv("DEV") == "YES" {
		err := Handler()
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	lambda.Start(Handler)
}
