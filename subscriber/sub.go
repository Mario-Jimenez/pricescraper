package subscriber

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Subscriber interface {
	Fetch(context.Context) (*Message, error)
	Commit(context.Context, *Message) error
	Close() error
}

type handler struct {
	subscriber Subscriber
	fn         func([]byte)
}

func NewHandler(subscriber Subscriber, fn func([]byte)) *handler {
	return &handler{subscriber, fn}
}

// InboundMessages handles inbound messages from the broker
func (h *handler) InboundMessages() {
	ctx := context.Background()
	fwait := 1
	for {
		m, err := h.subscriber.Fetch(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"wait":  fmt.Sprintf("Retrying in %d second(s)", fwait),
				"error": err.Error(),
			}).Warning("Failed to fetch message from the broker. Retrying...")
			time.Sleep(time.Duration(fwait) * time.Second)
			if fwait <= 60 {
				fwait += 3
			}
			continue
		}

		fwait = 1

		h.fn(m.Message)

		cwait := 1
		for {
			if err := h.subscriber.Commit(ctx, m); err != nil {
				log.WithFields(log.Fields{
					"wait":  fmt.Sprintf("Retrying in %d second(s)", cwait),
					"error": err.Error(),
				}).Warning("Failed to commit message's offset to the broker. Retrying...")
				time.Sleep(time.Duration(cwait) * time.Second)
				if cwait <= 60 {
					cwait += 3
				}
				continue
			}
			break
		}
	}
}
