package register

import "sync"

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]ForexService)
)

func Register(name string, driver ForexService) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("forex: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("forex: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(name string) ForexService {
	return drivers[name]
}

type Currency string

const (
	USD Currency = "USD"
	CAD Currency = "CAD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// ForexService is an interface that provides methods for querying foreign exchange rates.
type ForexService interface {
	// GetRate returns the exchange rate of a currency pair in basis points or an error if one occurred.
	GetRate(c1, c2 Currency) (int, error)
}
