package logger

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func BackgroundInfo(msg ...interface{}) {
	log.Info().Msg(fmt.Sprint(msg...))
}

func BackgroundInfof(msg string, parameters ...interface{}) {
	log.Info().Msg(fmt.Sprintf(msg, parameters...))
}

func BackgroundError(err error, msg ...interface{}) {
	log.Error().Err(err).Msg(fmt.Sprint(msg...))
}

func BackgroundErrorf(err error, msg string, parameters ...interface{}) {
	log.Error().Err(err).Msg(fmt.Sprintf(msg, parameters...))
}

func Info(ctx context.Context, msg ...interface{}) {
	log.Info().Str("TraceID", getTraceId(ctx)).Msg(fmt.Sprint(msg...))
}

func Infof(ctx context.Context, msg string, parameters ...interface{}) {
	log.Info().Str("TraceID", getTraceId(ctx)).Msg(fmt.Sprintf(msg, parameters...))
}

func Error(ctx context.Context, err error, msg ...interface{}) {
	log.Error().Str("TraceID", getTraceId(ctx)).Err(err).Msg(fmt.Sprint(msg...))
}

func Errorf(ctx context.Context, err error, msg string, parameters ...interface{}) {
	log.Error().Str("TraceID", getTraceId(ctx)).Err(err).Msg(fmt.Sprintf(msg, parameters...))
}

func Errors(ctx context.Context, err []error, msg ...interface{}) {
	log.Error().Str("TraceID", getTraceId(ctx)).Errs("errors", err).Msg(fmt.Sprint(msg...))
}

func Errorsf(ctx context.Context, err []error, msg string, parameters ...interface{}) {
	log.Error().Str("TraceID", getTraceId(ctx)).Errs("errors", err).Msg(fmt.Sprintf(msg, parameters...))
}

func getTraceId(ctx context.Context) string {
	return fmt.Sprint(ctx.Value("traceId"))
}
