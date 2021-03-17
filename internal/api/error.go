package api

import (
	"github.com/aws/aws-lambda-go/events"

	"log"
)

func ApiGatewayProxyResponseError(status int, err error) (events.APIGatewayProxyResponse, error) {
	log.Printf("Error %d: %s", status, err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, err
}
