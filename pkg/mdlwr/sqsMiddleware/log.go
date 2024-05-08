package sqsmiddleware

import (
	"context"
	"main/pkg/mdlwr"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func RequestLog(next mdlwr.SQSHandlerFunction) mdlwr.SQSHandlerFunction {
	return func(ctx context.Context, msg events.SQSEvent) error {
		message := msg.Records[0]
		log.Info().Str("trace_id", mdlwr.GetTraceId(ctx)).
			Msg("Processing SQS message with body: " + message.Body)

		h := next(ctx, msg)
		return h
	}
}
