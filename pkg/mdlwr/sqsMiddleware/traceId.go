package sqsmiddleware

import (
	"context"
	"main/pkg/mdlwr"

	"github.com/aws/aws-lambda-go/events"
	"github.com/segmentio/ksuid"
)

func TraceId(next mdlwr.SQSHandlerFunction) mdlwr.SQSHandlerFunction {
	return func(ctx context.Context, msg events.SQSEvent) error {
		var traceId string
		message := msg.Records[0]
		if _, ok := message.MessageAttributes["traceId"]; !ok {
			traceId = ksuid.New().String()
		} else {
			traceId = *message.MessageAttributes["traceId"].StringValue
		}
		ctx = context.WithValue(ctx, "traceId", traceId)
		return next(ctx, msg)
	}
}
