package scraper

type (
	Game struct {
		ID         string  `json:"id"`
		BasePrice  float32 `json:"base_price"`
		PricesURLs []price `json:"prices_urls"`
	}

	price struct {
		From string `json:"from"`
		URL  string `json:"url"`
	}
)
