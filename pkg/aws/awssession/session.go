package awssession

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	URL       string
}

var awsSession *session.Session

func Initialize(conf Config) {
	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			conf.AccessKey,
			conf.SecretKey,
			"",
		),
		Region: &conf.Region,
	}
	if conf.URL != "" {
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	awsSession = session.Must(session.NewSession(awsConfig))
}

func GetSession() *session.Session {
	return awsSession
}
