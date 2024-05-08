package mdlwr

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error)
type SQSHandlerFunction func(ctx context.Context, record events.SQSEvent) error
type MiddlewareFunc func(next HandlerFunc) HandlerFunc
type MiddlewareSQSFunc func(next SQSHandlerFunction) SQSHandlerFunction
