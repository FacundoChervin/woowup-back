package emailsrepo

import (
	"context"
	"fmt"
	"main/internal/core/domain"
	"main/internal/core/ports"
	"main/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func CreateRepository(sess *session.Session, dynamo dynamodbiface.DynamoDBAPI, tableName string) ports.EmailRepository {
	return emailRepository{
		TableName: tableName,
		dynamodb:  dynamo,
		session:   sess,
	}
}

type emailRepository struct {
	TableName string
	dynamodb  dynamodbiface.DynamoDBAPI
	session   *session.Session
}

func (r emailRepository) Get(ctx context.Context, id string) (*[]domain.EmailEntity, error) {
	keyExpr := expression.Key("batchId").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpr).Build()
	if err != nil {
		return nil, err
	}
	query, err := r.dynamodb.QueryWithContext(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	returnEmails := []domain.EmailEntity{}

	for _, i := range query.Items {
		var item domain.EmailEntity
		err := dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			logger.Errorf(ctx, err, "Error unmarshalling item:")
			return nil, err
		}
		returnEmails = append(returnEmails, item)
	}
	return &returnEmails, nil
}

func (r emailRepository) Save(ctx context.Context, email domain.EmailEntity) error {
	av, err := dynamodbattribute.MarshalMap(email)
	if err != nil {
		return fmt.Errorf("got error marshalling new email item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.TableName),
	}

	_, err = r.dynamodb.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %s", err)
	}

	return nil
}

func (r emailRepository) Find(ctx context.Context, partitionKey string, sortKey string) (*domain.EmailEntity, error) {
	av, err := dynamodbattribute.MarshalMap(EmailKey{PartitionKey: &partitionKey, SortKey: &sortKey})
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       av,
		TableName: aws.String(r.TableName),
	}

	result, err := r.dynamodb.GetItem(input)
	if err != nil {
		return nil, err
	}
	var response domain.EmailEntity
	dynamodbattribute.UnmarshalMap(result.Item, &response)
	return &response, nil
}
