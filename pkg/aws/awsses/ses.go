package awsses

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var sesInstance *ses.SES

func Initialize(awsSession *session.Session) {
	sesInstance = ses.New(awsSession)
}

func GetInstance() *ses.SES {
	return sesInstance
}
