package main

import (
	"main/internal/connectors"
	"main/internal/connectors/awssesconnector"
	"main/internal/connectors/mailgun"
	"main/internal/connectors/sendgrid"
	emailssrvc "main/internal/core/services/emailssrv"
	"main/internal/core/services/messagesrv"
	"main/internal/handlers/emailssqs"
	"main/internal/repositories/emailsrepo"
	"main/pkg/aws/awsses"
	"main/pkg/aws/awssession"
	"main/pkg/aws/dynadb"
	"main/pkg/mdlwr"
	sqsmiddleware "main/pkg/mdlwr/sqsMiddleware"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spf13/viper"
)

var handler mdlwr.SQSHandlerFunction

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
	emailsProviderList := make(map[string]connectors.EmailProviderService)
	emailsProviderList["ses"] = awssesconnector.CreateSESService(awsses.GetInstance())
	emailsProviderList["mailgun"] = mailgun.CreateMailgunService(viper.GetString("MAILGUN_API_KEY"), viper.GetString("EMAIL_DOMAIN"))
	emailsProviderList["sendgrid"] = sendgrid.CreateSendgridService(viper.GetString("SENDGRID_API_KEY"))

	emailProviderEnv := viper.GetString("EMAIL_PROVIDERS_ORDER")
	emailProvidersOrder := strings.Split(emailProviderEnv, ",")

	for _, v := range emailProvidersOrder {
		emailsProviders = append(emailsProviders, emailsProviderList[v])
	}

	sqsQueueName := viper.GetString("SEND_SQS_QUEUE_NAME")
	messageService := messagesrv.CreateSenderService(sqsQueueName)

	emailsrepo := emailsrepo.CreateRepository(awsSession, dynadb.GetInstance(), viper.GetString("LAMBDA_AWS_DYNAMO_TABLE"))
	emailssrv := emailssrvc.CreateService(emailsrepo, emailsProviders, messageService)
	emailshdl := emailssqs.CreateSQSController(emailssrv)

	mw := mdlwr.CreateLambdaSQSMiddlewareBuilder()
	mw.AddMiddleware(sqsmiddleware.TraceId)
	mw.AddMiddleware(sqsmiddleware.RequestLog)
	handler = mw.Build(emailshdl.Send)

}

func main() {
	lambda.Start(handler)
}
