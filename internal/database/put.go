package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"

	"fmt"
	"time"
)

const (
	DateFormat string = "20060102"
)

func Connect(input ConnectionInput) (conn Connection, err error) {

	conn.EtlResultsTable = input.EtlResultsTable
	conn.YahooTrendingTable = input.YahooTrendingTable
	conn.RedditResponseArchiveTable = input.RedditResponseArchiveTable

	conn.Session = input.Session
	conn.Service = dynamodb.New(conn.Session)

	return conn, err
}

func NewRedditResponseArchiveTable(p []reddit.Post, sub string, t time.Time) (r RedditResposeArchiveRecord) {

	for _, v := range p {
		r.Posts = append(r.Posts, v.Permalink)
	}
	r.Id = fmt.Sprintf("%s_%s", sub, t.Format(DateFormat))
	r.Hour = t.Hour()

	return r
}

func (c *Connection) PutRedditResonseArchiveRecord(record RedditResposeArchiveRecord) (err error) {

	if c.EtlResultsTable == "" {
		return fmt.Errorf("RedditResponseArchiveTable name is undefined")
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.RedditResponseArchiveTable),
	}

	_, err = c.Service.PutItem(input)

	return err
}

func NewEtlResultsRecord(sr []report.StockReport, sub string, t time.Time) (e EtlResultsRecord) {

	e.Data = sr
	e.Id = fmt.Sprintf("%s_%s", sub, t.Format(DateFormat))
	e.Hour = t.Hour()

	return e
}

func (c *Connection) PutEtlResultsRecord(record EtlResultsRecord) (err error) {

	if c.EtlResultsTable == "" {
		return fmt.Errorf("EtlResultsTable name is undefined")
	}

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

	if c.EtlResultsTable == "" {
		return record, fmt.Errorf("EtlResultsTable name is undefined")
	}

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
