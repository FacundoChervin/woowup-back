package awssesconnector

import (
	"bytes"
	"log"
	"main/internal/connectors"
	"main/internal/core/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

func CreateSESService(sesSession *ses.SES) connectors.EmailProviderService {
	return sesService{
		sesSession: sesSession,
	}
}

type sesService struct {
	sesSession *ses.SES
}

func (s sesService) SendEmail(emailData domain.SendEmail) (domain.EmailView, error) {

	msg := gomail.NewMessage()

	msg.SetAddressHeader("From", emailData.Email.From.Email, emailData.Email.From.Name)
	msg.SetHeader("To", []string{(emailData.Recipient.Email)}...)
	msg.SetHeader("Subject", emailData.Email.Subject)
	msg.SetBody("text/html", emailData.Email.HtmlContent)

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	message := ses.RawMessage{Data: emailRaw.Bytes()}
	input := &ses.SendRawEmailInput{Source: aws.String(emailData.Email.From.Email), Destinations: []*string{aws.String(emailData.Recipient.Email)}, RawMessage: &message}

	_, err := s.sesSession.SendRawEmail(input)
	if err != nil {
		log.Println(err)
		return domain.EmailView{}, err
	}

	return domain.EmailView{ID: emailData.Email.ID}, nil
}
