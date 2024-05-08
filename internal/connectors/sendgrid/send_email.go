package sendgrid

import (
	"fmt"
	"log"
	"main/internal/connectors"
	"main/internal/core/domain"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func CreateSendgridService(apiKey string) connectors.EmailProviderService {
	return sendGridSevice{
		apiKey: apiKey,
	}
}

type sendGridSevice struct {
	apiKey string
}

func (s sendGridSevice) SendEmail(emailData domain.SendEmail) (domain.EmailView, error) {
	from := mail.NewEmail(emailData.Email.From.Name, emailData.Email.From.Email)

	to := mail.NewEmail(emailData.Recipient.Name, emailData.Recipient.Email)
	message := mail.NewSingleEmail(from, emailData.Email.Subject, to, emailData.Email.PlainContent, emailData.Email.HtmlContent)

	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return domain.EmailView{}, err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return domain.EmailView{ID: emailData.Email.ID}, nil
}
