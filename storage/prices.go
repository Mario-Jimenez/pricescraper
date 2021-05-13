package storage

import (
	"sync"
)

type (
	Prices struct {
		// concurrency safe
		mu    sync.Mutex
		price float32
		url   string
	}
)

func NewPrices() *Prices {
	return &Prices{}
}

func (p *Prices) Add(url string, newPrice float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.price > 0 && p.price > newPrice && newPrice > 0 {
		p.price = newPrice
		p.url = url
		return
	}

	if p.price == 0 {
		p.price = newPrice
		p.url = url
	}
}

func (p *Prices) Get() (float32, string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.price, p.url
}

func (p *Prices) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.price = 0
	p.url = ""
}
