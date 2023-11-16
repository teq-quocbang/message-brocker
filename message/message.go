package message

import "context"

type MessageSender interface {
	Notification(ctx context.Context, msg string) error
}

type MessageReceiver interface {
}
