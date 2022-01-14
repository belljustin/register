package currencyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/belljustin/register"
)

func init() {
	register.Register("currencyapi", &ForexService{})
}

type ForexService struct {
}

func (s *ForexService) GetRate(c1, c2 register.Currency) (int, error) {
	if c1 == c2 {
		return 1_000_000, nil
	}

	lower_c1, lower_c2 := strings.ToLower(string(c1)), strings.ToLower(string(c2))
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/%s/%s.json", lower_c1, lower_c2)
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return -1, err
	}

	rate := body[lower_c2].(float64)
	return int(rate * 1_000_000), nil
}
