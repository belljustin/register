package main

import (
	"fmt"
	"strings"

	"github.com/belljustin/register"
	_ "github.com/belljustin/register/drivers/currencyapi"
	_ "github.com/belljustin/register/drivers/freeforex"
)

func main() {
	forexService := register.Open("freeforex")

	for {
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			panic(err)
		}
		currencies := strings.Split(input, ",")
		if len(currencies) != 2 {
			panic("requires two currency codes as input")
		}
		c1, c2 := register.Currency(currencies[0]), register.Currency(currencies[1])

		rate, err := forexService.GetRate(c1, c2)
		if err != nil {
			panic(err)
		}
		fmt.Println(rate)
	}
}
