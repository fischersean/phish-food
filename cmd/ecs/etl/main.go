package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/stocks"

	"log"
	"os"
	"sync"
	"time"
)

type subProcessingInput struct {
	Sub              string
	TickerPopulation []stocks.Stock
	AuthToken        reddit.AuthToken
	DBConnection     database.Connection
	PostLimit        int
	Wg               *sync.WaitGroup
}

var (
	TargetSubs []string = []string{
		"stocks",
		"wallstreetbets",
		"investing",
		"Wallstreetbetsnew",
		"WallStreetbetsELITE",
	}
	StartTime time.Time = time.Now()
)

const (
	PostCount int = 10
)

func processSub(input subProcessingInput) {

	defer input.Wg.Done()

	log.Printf("Downloading posts from %s", input.Sub)
	posts, err := getPosts(input.Sub, input.PostLimit, input.AuthToken, input.DBConnection)
	if err != nil {
		log.Printf("Could not fetch records for sub %s: %s", input.Sub, err.Error())
		return
	}

	log.Printf("Counting posts from %s", input.Sub)
	report, err := countRef(posts, input.TickerPopulation)
	if err != nil {
		log.Printf("Could not count post ref for %s: %s", input.Sub, err.Error())
		return
	}

	if os.Getenv("DEV") == "YES" {
		log.Printf("%s: %#v\n\n\n", input.Sub, report)
		return
	}

	log.Printf("Storing results from %s", input.Sub)
	err = putRecord(input.DBConnection, report, input.Sub, time.Now())
	if err != nil {
		log.Printf("Could not update database for sub %s: %s", input.Sub, err.Error())
		return
	}
}

func main() {

	bucketName := os.Getenv("BUCKET")
	appId := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	etlResultsTable := os.Getenv("TABLE")
	redditResponseArchveTable := os.Getenv("ARCHIVE_TABLE")

	// Add shared session to context
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	S3Service := s3.New(sess)
	tickerPopulation, err := getTradeableSecurities(S3Service, bucketName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Connect to db
	conn, err := database.Connect(database.ConnectionInput{
		Session:                    sess,
		EtlResultsTable:            etlResultsTable,
		RedditResponseArchiveTable: redditResponseArchveTable,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	auth, err := reddit.GetAuthToken(appId, appSecret)
	if err != nil {
		log.Fatal(err.Error())
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(TargetSubs))

	for _, sub := range TargetSubs {
		input := subProcessingInput{
			Wg:               wg,
			Sub:              sub,
			AuthToken:        auth,
			DBConnection:     conn,
			PostLimit:        PostCount,
			TickerPopulation: tickerPopulation,
		}
		go processSub(input)
	}

	wg.Wait()
}
