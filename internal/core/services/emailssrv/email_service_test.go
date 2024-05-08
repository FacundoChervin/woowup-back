package emailssrvc

import (
	"context"
	"main/internal/connectors"
	"main/internal/core/domain"
	"main/mocks"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

// Get email batches found
func Test_Get_Found(t *testing.T) {
	emailsRepo := mocks.CreateEmailRepository()

	emailsProviders := []connectors.EmailProviderService{}

	emailSvc := mocks.CreateEmailService(emailsRepo, emailsProviders)
	res, err := emailSvc.Get(context.TODO(), mocks.BatchIdFound)
	assert.NoError(t, err)
	for _, v := range *res {
		assert.Equal(t, v.ID, mocks.BatchIdFound)
	}
}

// Get email batches not found, no error
func Test_Get_Not_Found(t *testing.T) {
	emailsRepo := mocks.CreateEmailRepository()
	emailsProviders := []connectors.EmailProviderService{}

	emailSvc := mocks.CreateEmailService(emailsRepo, emailsProviders)
	res, err := emailSvc.Get(context.TODO(), ksuid.New().String())
	assert.NoError(t, err)
	assert.Nil(t, res)
}

// Email Sent no error
func Test_Sent_No_Error(t *testing.T) {
	emailsRepo := mocks.CreateEmailRepository()
	emailsProviders := []connectors.EmailProviderService{}

	emailsProviders = append(emailsProviders, mocks.CreateDefaultEmailProviderMockService(mocks.ProviderFail_1))
	emailsProviders = append(emailsProviders, mocks.CreateDefaultEmailProviderMockService(mocks.ProviderFail_2))
	emailsProviders = append(emailsProviders, mocks.CreateDefaultEmailProviderMockService(mocks.ProviderSuccess))

	emailSvc := CreateService(emailsRepo, emailsProviders, mocks.CreateSenderMockService(mocks.SqsQueueName))
	err := emailSvc.Send(context.TODO(), domain.SendEmail{})
	assert.NoError(t, err)
}
