package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"crypto/sha256"
	"fmt"
)

// GetKeyPermissions checks whether an *unhashed* key is valid and has the correct permissions to access the given route
func (c *Connection) GetKeyPermissions(input ApiKeyQueryInput) (record ApiKeyRecord, err error) {

	if input.UnhashedKey == "" {
		return record, fmt.Errorf("No key given in input")
	}

	h := sha256.New()
	_, err = h.Write([]byte(input.UnhashedKey))
	if err != nil {
		return record, err
	}

	qInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(fmt.Sprintf("%x", h.Sum(nil))),
			},
		},
		KeyConditionExpression: aws.String("key_hash = :v1"),
		TableName:              aws.String(c.ApiKeyTable),
	}

	result, err := c.Service.Query(qInput)
	if err != nil {
		return record, err
	}

	if *result.Count < 1 {
		return record, fmt.Errorf("Could not find key in database")
	}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &record)

	return record, err
}
