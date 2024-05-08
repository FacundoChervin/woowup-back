package dynadb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamoInstance *dynamodb.DynamoDB

func Initialize(awsSession *session.Session) {
	dynamoInstance = dynamodb.New(awsSession)
}

func GetInstance() *dynamodb.DynamoDB {
	return dynamoInstance
}
