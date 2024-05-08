package emailssqs

import (
	"context"
	"errors"
	"main/internal/core/domain"
	"main/internal/core/ports"
	"main/pkg/lambda/lambdasqs"
	"main/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/dealancer/validate.v2"
)

func CreateSQSController(emailService ports.EmailService) emailHandler {
	return emailHandler{
		emailService: emailService,
	}
}

type emailHandler struct {
	emailService ports.EmailService
}

// Processes the queue and tries to send the emails
func (gh emailHandler) Send(ctx context.Context, msg events.SQSEvent) error {
	var err error
	if len(msg.Records) < 1 {
		return errors.New("message batch empty")
	}

	message := msg.Records[0]
	requestData := new(sendEmailSqs)
	if err := lambdasqs.BindData(message, requestData); err != nil {
		return errors.New(bodyFormatErrorText)
	}

	if err := validate.Validate(requestData); err != nil {
		return err
	}
	emailData := domain.SendEmail{
		Email:     *requestData.Email,
		Provider:  *requestData.Provider,
		Recipient: *requestData.Destination,
	}

	err = gh.emailService.Send(ctx, emailData)
	logger.Error(ctx, err)

	return err

}
