package mdlwr

import (
	"context"
	"errors"
	"main/pkg/lambda"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/spf13/viper"
)

func AuthMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
		authToken, ok := request.Headers["authorization"]
		if !ok || authToken != ("Bearer "+viper.GetString("API_KEY")) {
			return ErrAuthMiddleware(ctx, request)
		}
		return next(ctx, request)
	}
}

func ErrAuthMiddleware(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	return lambda.LambdaErrorResponse(ctx, errors.New("unauthorized"), http.StatusForbidden), nil
}
