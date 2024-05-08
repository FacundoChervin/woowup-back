package messagesrv

import "context"

type SQSSender interface {
	SendToSQS(ctx context.Context, emailData any) error
}
