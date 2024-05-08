package mocks

import (
	"context"
	"fmt"
	"main/internal/core/domain"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func CreateDynamoClient() DynamodbMockClient {
	return DynamodbMockClient{}
}

type DynamodbMockClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m DynamodbMockClient) QueryWithContext(context.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error) {

	BatchIdFound := "BatchIdFound"
	data := []domain.EmailEntity{{
		ID:        BatchIdFound,
		From:      EmailFoundSent.From,
		Subject:   EmailFoundSent.Subject,
		Recipient: EmailFoundSent.Recipient,
		Sent:      EmailFoundSent.Sent,
	}, {
		ID:        BatchIdFound,
		From:      EmailFoundNotSent.From,
		Subject:   EmailFoundNotSent.Subject,
		Recipient: EmailFoundNotSent.Recipient,
		Sent:      EmailFoundNotSent.Sent,
	}}
	array := []map[string]*dynamodb.AttributeValue{}
	for _, v := range data {
		av, err := dynamodbattribute.MarshalMap(v)
		if err != nil {
			panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		}

		array = append(array, av)
	}

	return &dynamodb.QueryOutput{
		Items: array,
	}, nil
}
