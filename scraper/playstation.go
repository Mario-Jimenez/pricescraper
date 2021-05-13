package scraper

import (
	"strconv"
	"sync"

	"github.com/Mario-Jimenez/pricescraper/storage"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type PlayStationScraper struct {
	collector *colly.Collector
	prices    *storage.Prices
}

func NewPlayStationScraper(collector *colly.Collector, prices *storage.Prices) *PlayStationScraper {
	scraper := &PlayStationScraper{collector, prices}

	// Callback when colly finds the entry point to the DOM segment having a game info
	scraper.collector.OnHTML(`#main span.psw-h3`, func(e *colly.HTMLElement) {
		stringPrice := rex.ReplaceAllString(e.Text, "")
		if len(stringPrice) == 0 {
			stringPrice = "0"
		}
		gamePrice, err := strconv.ParseFloat(stringPrice, 32)
		if err != nil {
			log.WithFields(log.Fields{
				"text":  e.Text,
				"error": err.Error(),
			}).Error("Failed to parse float")
		}
		scraper.prices.Add(e.Request.URL.String(), float32(gamePrice))
	})

	scraper.collector.OnError(func(e *colly.Response, err error) {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("collector failed")
	})

	return scraper
}

func (s *PlayStationScraper) FindPrice(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	// start scraping the page under the given URL
	s.collector.Visit(url)
}
