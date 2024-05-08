package mocks

import (
	"errors"
	"main/internal/connectors"
	"main/internal/core/domain"
)

var (
	ProviderFail_1  = 0
	ProviderFail_2  = 1
	ProviderSuccess = 2
)

func CreateDefaultEmailProviderMockService(providerThatWorks int) connectors.EmailProviderService {
	return emailProviderMock{
		providerThatWorks: providerThatWorks,
	}
}

type emailProviderMock struct {
	providerThatWorks int
}

func (s emailProviderMock) SendEmail(emailData domain.SendEmail) (domain.EmailView, error) {
	if s.providerThatWorks == emailData.Provider {
		return domain.EmailView{}, nil
	}
	return domain.EmailView{}, errors.New("provider failed")
}
