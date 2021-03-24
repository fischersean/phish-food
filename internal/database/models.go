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
	RedditResponseArchiveTable string
	ApiKeyTable                string

	// YahooTrendingTable is deprecated
	YahooTrendingTable string

	// UserTable is deprecated
	UserTable string
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

type RedditResposeArchiveRecord struct {
	Id    string   `json:"id"`
	Hour  int      `json:"hour"`
	Posts []string `json:"data"`
}

type ConnectionInput struct {
	Session *session.Session

	// All the fields below are currently ignored
	EtlResultsTable            string
	YahooTrendingTable         string
	RedditResponseArchiveTable string
	UserTable                  string
}

type ApiKeyQueryInput struct {
	UnhashedKey string
}

type ApiKeyRecord struct {
	KeyHash     string   `json:"key_hash"`
	Permissions []string `json:"permissions"`
	Enabled     bool     `json:"enabled"`
}
type UserRecord struct {
	Username      string `json:"id"`
	ApiKey        string `json:"Key"`
	ApiKeyEnabled bool   `json:"KeyEnabled"`
}

type UserQueryInput struct {
	Username string
}

type UserUpdateInput struct {
	Username      string
	NewUserRecord UserRecord
}
