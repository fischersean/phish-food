package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/fischersean/phish-food/internal/report"

	"time"
)

type Connection struct {
	// Session is the shared aws session
	Session *session.Session

	// Service is the shared dynamo db connection
	Service *dynamodb.DynamoDB

	EtlResultsTable            string
	YahooTrendingTable         string
	RedditResponseArchiveTable string
}

type EtlResultsRecord struct {
	Id   string               `json:"id"`
	Hour int                  `json:"hour"`
	Data []report.StockReport `json:"data"`
}

type EtlResultsQueryInput struct {
	Subreddit string
	Date      time.Time
	Limit     int
}

type RedditResposeArchiveRecord struct {
	Id    string   `json:"id"`
	Hour  int      `json:"hour"`
	Posts []string `json:"data"`
}

type ConnectionInput struct {
	Session                    *session.Session
	EtlResultsTable            string
	YahooTrendingTable         string
	RedditResponseArchiveTable string
}
