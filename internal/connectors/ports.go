package connectors

import "main/internal/core/domain"

type EmailProviderService interface {
	SendEmail(emailData domain.SendEmail) (domain.EmailView, error)
}
