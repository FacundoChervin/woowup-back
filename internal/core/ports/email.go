package ports

import (
	"context"
	"main/internal/core/domain"
)

type EmailRepository interface {
	Save(ctx context.Context, email domain.EmailEntity) error
	Get(ctx context.Context, id string) (*[]domain.EmailEntity, error)
	Find(ctx context.Context, partitionKey string, sortKey string) (*domain.EmailEntity, error)
}

type EmailService interface {
	Send(ctx context.Context, car domain.SendEmail) error
	Get(ctx context.Context, id string) (*[]domain.EmailEntity, error)
	SendEmailsSQS(ctx context.Context, emailData domain.EmailBatch) error
}
