package main

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fischersean/phish-food/internal/database"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func validateRequest(request events.APIGatewayProxyRequest) error {
	return nil
}

func apiError(status int, err error) (events.APIGatewayProxyResponse, error) {
	log.Printf("Error %d: %s", status, err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, err
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	etlResultsTable := os.Getenv("TABLE")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	err := validateRequest(request)
	if err != nil {
		return apiError(400, err)
	}

	// TODO: make it possible to select a single hour
	subreddit := request.QueryStringParameters["subreddit"]
	date := request.QueryStringParameters["date"] // if not supplied, time.Now()

	var dateTime time.Time
	if date == "" {
		dateTime = time.Now()
	} else {
		dateTime, err = time.Parse(time.RFC3339, date)
		if err != nil {
			return apiError(400, err)
		}
	}

	conn, err := database.Connect(database.ConnectionInput{
		Session:         sess,
		EtlResultsTable: etlResultsTable,
	})
	if err != nil {
		return apiError(500, fmt.Errorf(""))
	}

	etlRecord, err := conn.GetEtlResultsRecordLatest(datEtlResultsLatestQueryInputyInput{
		Subreddit: subreddit,
		Date:      dateTime,
	})
	if err != nil {
		return apiError(400, err)
	}

	body, err := json.Marshal(etlRecord)
	if err != nil {
		return apiError(400, err)
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: buf.String(),
	}, err
}
func main() {
	lambda.Start(Handler)
}
