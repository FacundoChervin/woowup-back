package mocks

import (
	"context"
	"log"
)

var SqsQueueName = "test"

func CreateSenderMockService(queueName string) emailSenderService {
	if queueName == "" {
		log.Fatal("SQS queue name cannot be empty")
	}
	return emailSenderService{
		queueName: queueName,
	}
}

type emailSenderService struct {
	queueName string
}

func (s emailSenderService) SendToSQS(ctx context.Context, emailData any) error {
	return nil
}
