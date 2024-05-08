package emailsrepo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type EmailKey struct {
	PartitionKey *string `json:"batchId"`
	SortKey      *string `json:"recipient"`
}

func (pk *EmailKey) LoadFromDynamoKey(dynamoKey map[string]*dynamodb.AttributeValue) error {
	return dynamodbattribute.UnmarshalMap(dynamoKey, pk)
}
