package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"

	"fmt"
	//"strings"
	"bytes"
	"encoding/json"
	"time"
)

const (
	GetLatestRedditMaxLookback = 10
)

func NewRedditResponseArchiveRecord(p reddit.Post, sub string, t time.Time) (r RedditResposeArchiveRecord) {

	r.Post = p
	r.Permalink = p.Permalink
	r.Hour = t.Hour()
	r.Key = fmt.Sprintf("%s_%s/%d/%s.json", sub, t.Format(DateFormat), r.Hour, r.Permalink[:len(r.Permalink)-1])

	return r
}

func (c *Connection) PutRedditResponseArchiveRecord(record RedditResposeArchiveRecord) (err error) {

	buf, err := json.Marshal(record)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.RedditResponseArchiveBucket),
		Key:    aws.String(record.Key),
		Body:   bytes.NewReader(buf),
	}
	_, err = c.S3Service.PutObject(input)
	return err
}

func NewEtlResultsRecord(sr []report.StockReport, sub string, t time.Time) (e EtlResultsRecord) {

	e.Data = sr
	e.Id = fmt.Sprintf("%s_%s", sub, t.Format(DateFormat))
	e.Hour = t.Hour()

	return e
}

func (c *Connection) PutEtlResultsRecord(record EtlResultsRecord) (err error) {

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.EtlResultsTable),
	}

	_, err = c.Service.PutItem(input)

	return err
}

func (c *Connection) GetEtlResultsRecord(input EtlResultsQueryInput) (record []EtlResultsRecord, err error) {

	qInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(fmt.Sprintf("%s_%s", input.Subreddit, input.Date.Format(DateFormat))),
			},
		},
		KeyConditionExpression: aws.String("id = :v1"),
		TableName:              aws.String(c.EtlResultsTable),
	}

	result, err := c.Service.Query(qInput)

	if err != nil {
		return record, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &record)

	return record, err
}

func (c *Connection) GetLatestEtlResultsRecord(input EtlResultsQueryInput) (record EtlResultsRecord, err error) {

	// Starting from input.Date, loop backward by 1 day at a time until we get a result
	// Take the latest results based off of the sort key
	// If no date provided, start from time.Now()
	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	var result *dynamodb.QueryOutput
	count := int64(0)
	lookbackCount := 0

	for d := input.Date; count < 1; d.Add(-24 * time.Hour) {
		if lookbackCount > GetLatestRedditMaxLookback {
			return record, fmt.Errorf("Reached max lookback distance without finding primary key")
		}
		qInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v1": {
					S: aws.String(fmt.Sprintf("%s_%s", input.Subreddit, input.Date.Format(DateFormat))),
				},
			},
			KeyConditionExpression: aws.String("id = :v1"),
			ScanIndexForward:       aws.Bool(false),
			Limit:                  aws.Int64(1),
			TableName:              aws.String(c.EtlResultsTable),
		}

		result, err = c.Service.Query(qInput)
		if err != nil {
			return record, err
		}
		count = *result.Count
		lookbackCount += 1
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &record)

	return record, err
}
