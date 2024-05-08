package mdlwr

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/segmentio/ksuid"
)

func RequestIdMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
		id, ok := request.Headers["X-Request-Id"]
		if !ok {
			id = ksuid.New().String()
		}
		ctx = context.WithValue(ctx, "traceId", id)
		return next(ctx, request)
	}
}

func GetTraceId(ctx context.Context) string {
	return fmt.Sprint(ctx.Value("traceId"))
}
