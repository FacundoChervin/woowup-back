package mocks

import (
	"context"
	"fmt"
	"main/internal/connectors"
	"main/internal/core/domain"
)

func CreateEmailService(repository EmailMockRepository, emailProviders []connectors.EmailProviderService) EmailMockService {
	return EmailMockService{
		repository: repository,
	}
}

type EmailMockService struct {
	repository EmailMockRepository
}

func (s EmailMockService) Get(ctx context.Context, batchId string) (*[]domain.EmailEntity, error) {
	batch, err := s.repository.Get(ctx, batchId)
	if err != nil {
		return nil, fmt.Errorf("error from repository: %w", err)
	}
	return batch, nil
}

func (s EmailMockService) SendToSQS(ctx context.Context, emailData domain.SendEmail) error {
	return nil
}

func (s EmailMockService) Send(ctx context.Context, emailData domain.SendEmail) error {
	return nil
}
