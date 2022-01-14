package currencyapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/belljustin/register"
)

func TestGetRate(t *testing.T) {
	s := &ForexService{}
	rate, err := s.GetRate(register.USD, register.CAD)
	assert.Nil(t, err)
	assert.NotNil(t, rate)
}
