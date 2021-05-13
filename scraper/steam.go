package scraper

import (
	"strconv"
	"sync"

	"github.com/Mario-Jimenez/pricescraper/storage"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type SteamScraper struct {
	collector *colly.Collector
	prices    *storage.Prices
}

func NewSteamScraper(collector *colly.Collector, prices *storage.Prices) *SteamScraper {
	scraper := &SteamScraper{collector, prices}

	// Callback when colly finds the entry point to the DOM segment having a game info
	scraper.collector.OnHTML(`#game_area_purchase div.game_purchase_price.price`, func(e *colly.HTMLElement) {
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
			return
		}
		scraper.prices.Add(e.Request.URL.String(), float32(gamePrice))
	})

	// Callback when colly finds the entry point to the DOM segment having a game info
	scraper.collector.OnHTML(`#game_area_purchase div.discount_final_price`, func(e *colly.HTMLElement) {
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
			return
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

func (s *SteamScraper) FindPrice(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	// start scraping the page under the given URL
	s.collector.Visit(url)
}
