package publisher

import "context"

type Publisher interface {
	Publish(context.Context, []byte) error
	Close() error
}
