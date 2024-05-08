package main

import (
	"main/internal/connectors"
	emailssrvc "main/internal/core/services/emailssrv"
	"main/internal/core/services/messagesrv"
	"main/internal/handlers/emailshdl"
	"main/internal/repositories/emailsrepo"
	"main/pkg/aws/awsses"
	"main/pkg/aws/awssession"
	"main/pkg/aws/dynadb"
	"main/pkg/mdlwr"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spf13/viper"
)

var handler mdlwr.HandlerFunc

func init() {
	viper.AutomaticEnv()

	awsConfig := awssession.Config{
		AccessKey: viper.GetString("LAMBDA_AWS_ACCESS_KEY"),
		SecretKey: viper.GetString("LAMBDA_AWS_SECRET_KEY"),
		Region:    viper.GetString("LAMBDA_AWS_REGION"),
	}
	awssession.Initialize(awsConfig)
	awsSession := awssession.GetSession()
	dynadb.Initialize(awsSession)
	awsses.Initialize(awsSession)

	emailsProviders := []connectors.EmailProviderService{}

	sqsQueueName := viper.GetString("SEND_SQS_QUEUE_NAME")
	messageService := messagesrv.CreateSenderService(sqsQueueName)

	emailsrepo := emailsrepo.CreateRepository(awsSession, dynadb.GetInstance(), viper.GetString("LAMBDA_AWS_DYNAMO_TABLE"))
	emailssrv := emailssrvc.CreateService(emailsrepo, emailsProviders, messageService)
	emailshdl := emailshdl.CreateHTTPController(emailssrv)

	mw := mdlwr.CreateLambdaRESTMiddlewareBuilder()
	mw.AddMiddleware(mdlwr.AuthMiddleware)
	mw.AddMiddleware(mdlwr.RequestIdMiddleware)
	bodyExcludeKeys := []string{""}
	headerExcludeKeys := []string{""}
	logMiddleware := mdlwr.CreateLogMdlw(&bodyExcludeKeys, &headerExcludeKeys)
	mw.AddMiddleware(logMiddleware)

	handler = mw.Build(emailshdl.Send)
}

func main() {
	lambda.Start(handler)
}
