package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Mario-Jimenez/pricescraper/publisher"
	"github.com/Mario-Jimenez/pricescraper/storage"
	log "github.com/sirupsen/logrus"
)

var rex = regexp.MustCompile(`[^0-9\\.]+`)

type PriceScraper interface {
	FindPrice(string, *sync.WaitGroup)
}

type Handler struct {
	amazonScraper      PriceScraper
	nintendoScraper    PriceScraper
	steamScraper       PriceScraper
	playStationScraper PriceScraper
	prices             *storage.Prices
	publisher          publisher.Publisher
}

func NewHandler(amazonScraper, nintendoScraper, steamScraper, playStationScraper PriceScraper,
	prices *storage.Prices, publisher publisher.Publisher) *Handler {
	return &Handler{amazonScraper, nintendoScraper, steamScraper, playStationScraper, prices, publisher}
}

func (h *Handler) ProcessMessage(message []byte) {
	game := &Game{}
	err := json.Unmarshal(message, &game)
	if err != nil {
		log.WithFields(log.Fields{
			"message": string(message),
			"error":   err.Error(),
		}).Warning("An invalid message was received from the broker")
		return
	}

	log.WithFields(log.Fields{
		"game": game,
	}).Info("Inbound message from broker")

	wg := &sync.WaitGroup{}

	for _, g := range game.PricesURLs {
		wg.Add(1)
		switch g.From {
		case "amazon":
			go h.amazonScraper.FindPrice(g.URL, wg)
		case "steam":
			go h.steamScraper.FindPrice(g.URL, wg)
		case "playstation":
			go h.playStationScraper.FindPrice(g.URL, wg)
		case "nintendo":
			go h.nintendoScraper.FindPrice(g.URL, wg)
		default:
			log.WithFields(log.Fields{
				"from": g.From,
			}).Warning("invalid site")
			wg.Done()
		}
	}

	wg.Wait()

	finalPrice, finalURL := h.prices.Get()
	response := map[string]interface{}{
		"id":    game.ID,
		"price": finalPrice,
		"url":   finalURL,
	}

	// marshal prices
	b, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{
			"parameters": response,
			"error":      err.Error(),
		}).Error("json marshal failed")
	}

	h.prices.Reset()

	h.publishMessage(context.Background(), b)
}

func (h *Handler) publishMessage(ctx context.Context, message []byte) {
	// publish a message
	wait := 1
	for {
		if err := h.publisher.Publish(ctx, message); err != nil {
			log.WithFields(log.Fields{
				"message": string(message),
				"wait":    fmt.Sprintf("Retrying in %d second(s)", wait),
				"error":   err.Error(),
			}).Warning("Failed to publish message. Retrying...")
			time.Sleep(time.Duration(wait) * time.Second)
			if wait <= 60 {
				wait += 3
			}
			continue
		}

		log.WithFields(log.Fields{
			"message": string(message),
		}).Debug("Message published successfully")

		return
	}
}
