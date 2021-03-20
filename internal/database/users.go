package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
)

func (c *Connection) GetUserRecord(input UserQueryInput) (record UserRecord, err error) {

	if c.UserTable == "" {
		return record, fmt.Errorf("UserTable name is undefined")
	}

	qInput := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(input.Username),
			},
		},
		KeyConditionExpression: aws.String("id = :v1"),
		TableName:              aws.String(c.UserTable),
	}

	result, err := c.Service.Query(qInput)

	if err != nil {
		return record, err
	}

	if len(result.Items) == 0 {
		return record, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &record)

	return record, err
}

func (c *Connection) UpdateUserRecord(input UserUpdateInput) (record UserRecord, err error) {

	if c.UserTable == "" {
		return record, fmt.Errorf("UserTable name is undefined")
	}

	in := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#K": aws.String("Key"),
			"#E": aws.String("KeyEnabled"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":k": {
				S: aws.String(input.NewUserRecord.ApiKey),
			},
			":e": {
				BOOL: aws.Bool(input.NewUserRecord.ApiKeyEnabled),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(input.Username),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(c.UserTable),
		UpdateExpression: aws.String("SET #K = :k, #E = :e"),
	}

	result, err := c.Service.UpdateItem(in)
	if err != nil {
		return record, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &record)

	return record, err
}
