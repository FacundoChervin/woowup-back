package emailshdl

import (
	"context"
	"errors"
	"main/internal/core/domain"
	"main/internal/core/ports"
	"main/pkg/lambda"
	"main/pkg/lambda/lambdahttp"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/segmentio/ksuid"
	"gopkg.in/dealancer/validate.v2"
)

func CreateHTTPController(emailService ports.EmailService) emailHandler {
	return emailHandler{
		emailService: emailService,
	}
}

type emailHandler struct {
	emailService ports.EmailService
}

func (gh emailHandler) Send(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	requestData := new(sendEmail)

	if err := lambdahttp.BindData(request, requestData); err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusBadRequest), nil
	}

	if err := validate.Validate(requestData); err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusBadRequest), nil
	}
	batchID := ksuid.New().String()
	emailBatch := domain.EmailBatch{
		ID: batchID,
		Email: domain.Email{
			From:         *requestData.From,
			Subject:      requestData.Subject,
			PlainContent: requestData.PlainContent,
			HtmlContent:  requestData.HtmlContent,
		},
		Recipìents: *requestData.Destinations,
		Provider:   0,
	}

	err := gh.emailService.SendEmailsSQS(ctx, emailBatch)
	if err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusInternalServerError), nil
	}

	return lambda.LambdaSuccessResponse(ctx, http.StatusAccepted, lambda.ResponseData{Name: "batchID", Data: batchID}), nil
}

func (gh emailHandler) Get(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	requestData := new(getEmail)

	if err := lambdahttp.BindData(request, requestData); err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusBadRequest), nil
	}

	if err := validate.Validate(requestData); err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusBadRequest), nil
	}

	emailsData, err := gh.emailService.Get(ctx, requestData.ID)
	if err != nil {
		return lambda.LambdaErrorResponse(ctx, err, http.StatusInternalServerError), nil
	}
	if len(*emailsData) == 0 {
		return lambda.LambdaErrorResponse(ctx, errors.New("batch ID not found"), http.StatusNotFound), nil
	}

	return lambda.LambdaSuccessResponse(ctx, http.StatusOK, lambda.ResponseData{Name: "emails", Data: emailsData}), nil
}
