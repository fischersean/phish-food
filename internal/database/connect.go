package database

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"

	"fmt"
	"os"
)

const (
	DateFormat string = "20060102"
)

func Connect(input ConnectionInput) (conn Connection, err error) {

	conn.etlResultsTable = os.Getenv("ETL_RESULTS_TABLE")
	conn.redditPostArchiveBucket = os.Getenv("REDDIT_ARCHIVE_BUCKET")

	if conn.etlResultsTable == "" {
		return conn, fmt.Errorf("Could not find required table name ETL_RESULTS_TABLE")
	}

	if conn.redditPostArchiveBucket == "" {
		return conn, fmt.Errorf("Could not find required table name REDDIT_ARCHIVE_BUCKET")
	}

	conn.session = input.Session
	conn.service = dynamodb.New(conn.session)
	conn.s3Service = s3.New(conn.session)

	return conn, err
}
