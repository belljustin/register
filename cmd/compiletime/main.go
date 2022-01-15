package main

import (
	"fmt"
	"os"

	"github.com/belljustin/register"
	_ "github.com/belljustin/register/drivers/freeforex"
)

func main() {
	args := os.Args[1:]
	c1, c2 := register.Currency(args[0]), register.Currency(args[1])

	forexService := register.Open("freeforex")
	rate, err := forexService.GetRate(c1, c2)
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
