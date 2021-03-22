package database

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
	"os"
)

const (
	DateFormat string = "20060102"
)

func Connect(input ConnectionInput) (conn Connection, err error) {

	conn.EtlResultsTable = os.Getenv("ETL_RESULTS_TABLE")
	conn.RedditResponseArchiveTable = os.Getenv("REDDIT_ARCHIVE_TABLE")

	if conn.EtlResultsTable == "" {
		return conn, fmt.Errorf("Could not find required table name ETL_RESULTS_TABLE")
	}

	if conn.RedditResponseArchiveTable == "" {
		return conn, fmt.Errorf("Could not find required table name REDDIT_ARCHIVE_TABLE")
	}
	//conn.YahooTrendingTable = os.Getenv("YAHOO_TRENDING_TABLE")
	//conn.UserTable = os.Getenv("USER_TABLE")

	conn.Session = input.Session
	conn.Service = dynamodb.New(conn.Session)

	return conn, err
}
