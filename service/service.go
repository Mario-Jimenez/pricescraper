package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Mario-Jimenez/pricescraper/broker/kafka"
	"github.com/Mario-Jimenez/pricescraper/config"
	"github.com/Mario-Jimenez/pricescraper/logger"
	"github.com/Mario-Jimenez/pricescraper/scraper"
	"github.com/Mario-Jimenez/pricescraper/storage"
	"github.com/Mario-Jimenez/pricescraper/subscriber"
	"github.com/gocolly/colly"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// Run service
func Run(serviceName, serviceVersion string) {
	// load app configuration
	conf, err := config.NewFileConfig()
	if err != nil {
		if errors.IsNotFound(err) {
			log.WithFields(log.Fields{
				"error": errors.Details(err),
			}).Error("Configuration file not found")
			return
		}
		if errors.IsNotValid(err) {
			log.WithFields(log.Fields{
				"error": errors.Details(err),
			}).Error("Invalid configuration values")
			return
		}
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to retrieve secrets")
		return
	}

	// initialize logger
	logger.InitializeLogger(serviceName, serviceVersion, conf.Values().LogLevel)

	producer := kafka.NewProducer("gamesprices", conf.Values().KafkaConnection)
	consumer := kafka.NewConsumer("games", "prices", conf.Values().KafkaConnection)

	prices := storage.NewPrices()

	amazonScraper := scraper.NewAmazonScraper(colly.NewCollector(colly.AllowURLRevisit()), prices)
	nintendoScraper := scraper.NewNintendoScraper(colly.NewCollector(colly.AllowURLRevisit()), prices)
	steamScraper := scraper.NewSteamScraper(colly.NewCollector(colly.AllowURLRevisit()), prices)
	playStationScraper := scraper.NewPlayStationScraper(colly.NewCollector(colly.AllowURLRevisit()), prices)

	pricesScraper := scraper.NewHandler(amazonScraper, nintendoScraper, steamScraper, playStationScraper, prices, producer)

	sub := subscriber.NewHandler(consumer, pricesScraper.ProcessMessage)
	go sub.InboundMessages()

	log.Info("Running service...")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	log.Info("Shutting down server...")

	if err := consumer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close consumer")
	}

	if err := producer.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": errors.Details(err),
		}).Error("Failed to close publisher")
	}
}
