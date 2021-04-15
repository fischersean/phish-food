package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"

	"time"
)

type Connection struct {
	// session is the shared aws session
	session *session.Session

	// Service is the shared dynamo db connection
	service *dynamodb.DynamoDB

	// s3Service is the shared S3 connection
	s3Service *s3.S3

	etlResultsTable         string
	redditPostArchiveBucket string
}

type EtlResultsRecord struct {
	Id   string               `json:"id"`
	Hour int                  `json:"hour"`
	Data []report.StockReport `json:"data"`
}

type EtlResultsQueryInput struct {
	Subreddit string
	Date      time.Time
}

type RedditPostArchiveRecord struct {
	Key       string      `json:"key"`
	Hour      int         `json:"hour"`
	Permalink string      `json:"permalink,omitempty"`
	Post      reddit.Post `json:"data"`
}

type RedditPostArchiveListInput struct {
	Subreddit string
	Date      time.Time
}

type RedditPostArchiveQueryInput struct {
	Key string
}

type ConnectionInput struct {
	Session *session.Session

	// All the fields below are currently ignored
	EtlResultsTable        string
	YahooTrendingTable     string
	RedditPostArchiveTable string
	UserTable              string
}
