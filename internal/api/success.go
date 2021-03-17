package api

import (
	"github.com/aws/aws-lambda-go/events"
)

func ApiGatewayProxyResponsSuccess(body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: body,
	}, nil
}
