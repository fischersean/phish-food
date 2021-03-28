package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fischersean/phish-food/internal/etl"

	db "github.com/fischersean/phish-food/internal/database"
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
	DBConnection     db.Connection
	PostLimit        int
	Wg               *sync.WaitGroup
}

var (
	StartTime time.Time = time.Now()
)

const (
	PostCount        int = 100
	PostAnalyzeCount int = 10
)

func processSub(input subProcessingInput) {

	defer input.Wg.Done()

	log.Printf("Downloading posts from %s", input.Sub)
	posts, err := etl.GetPosts(input.Sub, input.PostLimit, input.AuthToken, input.DBConnection)
	if err != nil {
		log.Printf("Could not fetch records for sub %s: %s", input.Sub, err.Error())
		return
	}

	log.Printf("Counting posts from %s", input.Sub)
	var maxCount int
	if len(posts) < PostAnalyzeCount {
		maxCount = len(posts) - 1
	} else {
		maxCount = PostAnalyzeCount - 1
	}
	report, err := etl.CountRef(posts[0:maxCount], input.TickerPopulation)
	if err != nil {
		log.Printf("Could not count post ref for %s: %s", input.Sub, err.Error())
		return
	}

	if os.Getenv("DEV") == "YES" {
		log.Printf("%s: %#v\n\n\n", input.Sub, report)
		log.Printf("Dev environment detected. Skipping store step")
		return
	}

	log.Printf("Storing results from %s", input.Sub)
	err = etl.PutRecord(input.DBConnection, report, input.Sub, StartTime)
	if err != nil {
		log.Printf("Could not update database for sub %s: %s", input.Sub, err.Error())
		return
	}
}

func main() {

	bucketName := os.Getenv("TRADEABLES_BUCKET")
	appId := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")

	// Add shared session to context
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	S3Service := s3.New(sess)
	tickerPopulation, err := etl.GetTradeableSecurities(S3Service, bucketName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Connect to db
	conn, err := db.Connect(db.ConnectionInput{
		Session: sess,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	auth, err := reddit.GetAuthToken(appId, appSecret)
	if err != nil {
		log.Fatal(err.Error())
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(etl.FetchTargets))

	for _, sub := range etl.FetchTargets {
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
