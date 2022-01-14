package freeforex

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/belljustin/register"
)

var currencyPriorities = [...]register.Currency{register.EUR, register.GBP, register.USD, register.CAD}

func init() {
	register.Register("freeforex", &ForexService{})
}

type PairResponse struct {
	Rates map[string]struct {
		Rate float64
	}
}

type ForexService struct {
}

func (s *ForexService) GetRate(c1, c2 register.Currency) (int, error) {
	if c1 == c2 {
		return 1_000_000, nil
	}

	var c1Priority, c2Priority int
	for i, currency := range currencyPriorities {
		if c1 == currency {
			c1Priority = i
		} else if c2 == currency {
			c2Priority = i
		}
	}

	var pair string
	if c1Priority < c2Priority {
		pair = fmt.Sprintf("%s%s", c1, c2)
	} else {
		pair = fmt.Sprintf("%s%s", c2, c1)
	}

	url := fmt.Sprintf("https://www.freeforexapi.com/api/live?pairs=%s", pair)
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var pairResponse PairResponse
	if err := json.NewDecoder(resp.Body).Decode(&pairResponse); err != nil {
		return -1, err
	}

	rate, ok := pairResponse.Rates[pair]
	if !ok {
		return -1, fmt.Errorf("response did not include rate %s", pair)
	}

	var ret int
	if c1Priority < c2Priority {
		ret = int(rate.Rate * 1_000_000)
	} else {
		ret = int(math.Pow(rate.Rate, -1) * 1_000_000)
	}
	return ret, nil
}
