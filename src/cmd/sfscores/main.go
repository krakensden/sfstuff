package main

import (
	"fmt"
	"os"
	"sflib"
	"time"
)

func main() {
	client, stocks, err := sflib.Setup()
	if err != nil || client == nil || stocks == nil {
		fmt.Fprintln(os.Stderr, "setup failed", err)
		os.Exit(1)
	}
	target := os.Getenv("STOCK")
	account := os.Getenv("ACCOUNT")
	if target == "" || account == "" || err != nil {
		fmt.Println(os.Stderr, "Need STOCK and ACCOUNT environment variables to implement a buying strategy")
		os.Exit(2)
	}

	update := func(r chan sflib.StockOrders) {
		// Get Quote
		sq, err := client.CheckAllOrderStatus(*stocks.Venue, account, target)
		if err != nil || !sq.Ok {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprintln(os.Stderr, "Quote request error")
			}
		}
		r <- sq
	}
	results := make(chan sflib.StockOrders, 5)
	go update(results)
	fmt.Println("Starting the main loop")

	for {
		// check
		select {
		case <-time.After(1 * time.Second):
			go update(results)
		case r := <-results:
			var accounts map[string]int = make(map[string]int)
			for _, order := range r.Orders {
				currentCount, ok := accounts[order.Account]
				if !ok {
					accounts[order.Account] = 0
				}
				if order.Direction == sflib.Buy {
					accounts[order.Account] = currentCount + order.TotalFilled
				} else {
					accounts[order.Account] = currentCount - order.TotalFilled
				}
			}
			for account, filled := range accounts {
				fmt.Println("Account ", account, " owns ", filled, " shares")
			}
			fmt.Println("***************************************************")
		}
	}
}
