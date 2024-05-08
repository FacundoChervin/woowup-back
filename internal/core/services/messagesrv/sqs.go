package messagesrv

import (
	"context"
	"log"
	"main/pkg/lambda/lambdasqssender"
)

func CreateSenderService(queueName string) emailSenderService {
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
	err := lambdasqssender.Send(ctx, s.queueName, emailData)
	if err != nil {
		return err
	}
	return nil
}
