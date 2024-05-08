package emailssrvc

import (
	"context"
	"errors"
	"fmt"
	"main/internal/connectors"
	"main/internal/core/domain"
	"main/internal/core/ports"
	"main/internal/core/services/messagesrv"
	"main/pkg/logger"
)

func CreateService(repository ports.EmailRepository, emailProviders []connectors.EmailProviderService, sqsSender messagesrv.SQSSender) emailService {
	return emailService{
		repository:     repository,
		emailProviders: emailProviders,
		sqsSender:      sqsSender,
	}
}

type emailService struct {
	repository     ports.EmailRepository
	emailProviders []connectors.EmailProviderService
	sqsSender      messagesrv.SQSSender
}

func (s emailService) Send(ctx context.Context, emailData domain.SendEmail) error {
	var err error

	if len(s.emailProviders) <= 0 {
		return errors.New("no providers configured")
	}

	if emailData.Provider >= len(s.emailProviders) {
		return errors.New("invalid provider")
	}

	exist, err := s.repository.Find(ctx, emailData.Email.ID, emailData.Recipient.Email)
	if err != nil {
		logger.BackgroundErrorf(err, "error from repository when looking for email")
		return err
	}

	retries := 0
	if exist.ID != "" {
		retries = exist.Retries
	}

	emailToSave := domain.EmailEntity{
		ID:        emailData.Email.ID,
		Recipient: emailData.Recipient.Email,
		Content:   emailData.Email.PlainContent,
		From:      emailData.Email.From,
		Subject:   emailData.Email.Subject,
	}

	_, err = s.emailProviders[emailData.Provider].SendEmail(emailData)
	if err != nil {
		emailToSave.Retries = retries + 1
		emailToSave.Sent = false

		err = s.repository.Save(ctx, emailToSave)
		if err != nil {
			logger.BackgroundErrorf(err, "error from repository when saving email")
			return err
		}

		newEmailData := domain.SendEmail{
			Email:     emailData.Email,
			Provider:  (emailData.Provider + 1) % len(s.emailProviders),
			Recipient: emailData.Recipient,
		}
		err = s.sqsSender.SendToSQS(ctx, newEmailData)
		if err != nil {
			return err
		}
	} else {
		emailToSave.Sent = true
		emailToSave.Retries = retries
		err = s.repository.Save(ctx, emailToSave)
		if err != nil {
			logger.BackgroundErrorf(err, "error from repository when saving email")
			return err
		}
	}

	return nil
}

func (s emailService) SendEmailsSQS(ctx context.Context, emailData domain.EmailBatch) error {
	for _, v := range emailData.Recipìents {
		emailData.Email.ID = emailData.ID
		email := domain.SendEmail{
			Email:     emailData.Email,
			Provider:  emailData.Provider,
			Recipient: v,
		}
		err := s.sqsSender.SendToSQS(ctx, email)
		if err != nil {
			return fmt.Errorf("error when sending msg to create queue, %v", err)
		}
	}

	return nil
}

func (s emailService) Get(ctx context.Context, id string) (*[]domain.EmailEntity, error) {
	return s.repository.Get(ctx, id)
}
