package mailgun

import (
	"context"
	"fmt"
	"main/internal/connectors"
	"main/internal/core/domain"
	"main/pkg/logger"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func CreateMailgunService(apiKey string, domain string) connectors.EmailProviderService {
	return mailGunService{
		apiKey: apiKey,
		domain: domain,
	}
}

type mailGunService struct {
	apiKey string
	domain string
}

func (s mailGunService) SendEmail(emailData domain.SendEmail) (domain.EmailView, error) {
	mg := mailgun.NewMailgun(s.domain, s.apiKey)

	message := mg.NewMessage(emailData.Email.From.Email, emailData.Email.Subject, emailData.Email.PlainContent, emailData.Recipient.Email)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return domain.EmailView{}, err
	}
	logger.BackgroundInfo("SENDING EMAWIIIL")
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return domain.EmailView{ID: emailData.Email.ID}, nil
}
