package mocks

import (
	"context"
	"fmt"
	"main/internal/core/domain"
)

func CreateEmailRepository() EmailMockRepository {
	return EmailMockRepository{}
}

type EmailMockRepository struct {
}

var (
	BatchIdFound    = "2g6hxbfK2MCYqYDZiYm5Jo52Im6"
	RepositoryError = "error"
	EmailFoundSent  = domain.EmailEntity{
		From: domain.EmailAddress{
			Name:  "Facundo",
			Email: "facundo.chervin@gmail.com",
		},
		Recipient: "pedro.ramirez@gmail.com",
		Subject:   "Welcome message",
		Sent:      true,
	}
	EmailFoundNotSent = domain.EmailEntity{
		From: domain.EmailAddress{
			Name:  "Facundo",
			Email: "facundo.chervin@gmail.com",
		},
		Recipient: "juan.ramirez@gmail.com",
		Subject:   "Welcome message",
		Sent:      false,
	}
)

func (sr EmailMockRepository) Get(ctx context.Context, id string) (*[]domain.EmailEntity, error) {
	switch id {
	case BatchIdFound:
		return &[]domain.EmailEntity{{
			ID:        BatchIdFound,
			From:      EmailFoundSent.From,
			Subject:   EmailFoundSent.Subject,
			Recipient: EmailFoundSent.Recipient,
			Sent:      EmailFoundSent.Sent,
		}, {
			ID:        BatchIdFound,
			From:      EmailFoundNotSent.From,
			Subject:   EmailFoundNotSent.Subject,
			Recipient: EmailFoundNotSent.Recipient,
			Sent:      EmailFoundNotSent.Sent,
		}}, nil
	case RepositoryError:
		return nil, fmt.Errorf("error from repository: %w", fmt.Errorf("mock error"))
	default:
		return nil, nil
	}
}

func (sr EmailMockRepository) Save(ctx context.Context, email domain.EmailEntity) error {
	return nil
}

func (sr EmailMockRepository) Find(ctx context.Context, partitionKey string, sortKey string) (*domain.EmailEntity, error) {
	return &domain.EmailEntity{}, nil
}
