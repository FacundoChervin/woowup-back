package lambdasqssender

import (
	"context"
	"encoding/json"
	"fmt"
	"main/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/segmentio/ksuid"
)

type SQSConfig struct {
	MaxNumberOfMessages int64  `mapstructure:"SQS_MAX_NUMBER_OF_MESSAGES"`
	QueueName           string `mapstructure:"SQS_QUEUE_NAME"`
	QueueURL            string
	WaitTimeSeconds     int64 `mapstructure:"SQS_WAIT_TIME_SECONDS"`
}

type SQSSender struct {
	config    SQSConfig
	sqsClient *sqs.SQS
}

func Send(ctx context.Context, queueName string, body any) error {
	brokerSender := createSQSSender(queueName)
	return brokerSender.sendRawMessage(ctx, body)
}

func SendMultiple(ctx context.Context, queueName string, body any) error {
	brokerSender := createSQSSender(queueName)
	bodies, ok := body.([]any)
	if !ok {
		return fmt.Errorf("error trying to send multiple sqs messages, body must be an array")
	}
	for _, body := range bodies {
		err := brokerSender.sendRawMessage(ctx, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func getQueueURL(client *sqs.SQS, queueName string) (queueURL string) {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	response, err := client.GetQueueUrl(params)
	if err != nil {
		logger.BackgroundError(err, "Error when obtaining the queue url")
		return ""
	}
	queueURL = aws.StringValue(response.QueueUrl)

	return queueURL
}

func CreateSqsClient() *sqs.SQS {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sqs.New(awsSession, awsSession.Config)
}

func createSQSSender(queueName string) SQSSender {
	defaultSqsConf := SQSConfig{
		MaxNumberOfMessages: 10,
		QueueName:           queueName,
		WaitTimeSeconds:     5,
	}

	sqsSender := SQSSender{
		config: defaultSqsConf,
	}

	client := CreateSqsClient()
	sqsSender.config.QueueURL = getQueueURL(client, queueName)
	sqsSender.sqsClient = client
	return sqsSender
}

func (sqsSender SQSSender) sendRawMessage(ctx context.Context, body any) error {

	traceId := ctx.Value("traceId").(string)
	if traceId == "" {
		traceId = ksuid.New().String()
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		logger.BackgroundError(err, "Error when marshaling body")
		return err
	}

	msgAttributes := map[string]*sqs.MessageAttributeValue{
		"traceId": {
			DataType:    aws.String("String"),
			StringValue: &traceId,
		},
	}
	bodyStr := string(bodyByte)
	params := &sqs.SendMessageInput{
		MessageBody:       &bodyStr,
		QueueUrl:          &sqsSender.config.QueueURL,
		MessageAttributes: msgAttributes,
	}
	dupliId := ksuid.New().String()
	groupId := ksuid.New().String()

	params.MessageDeduplicationId = &dupliId
	params.MessageGroupId = &groupId

	_, er := sqsSender.sqsClient.SendMessage(params)

	if er == nil {
		logger.Infof(ctx, "Sent message to SQS: %v", params)
	} else {
		logger.Errorf(ctx, er, "Error when sending SQS message: %v", params)
	}

	return er
}
