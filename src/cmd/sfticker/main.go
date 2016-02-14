package main

import (
	"fmt"
	"os"
	"sflib"
	"time"
)

func main() {
	sfc, vs, err := sflib.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup failed", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Available stocks")
	for {
		// check
		results := make(chan sflib.StockQuote, len(vs.Symbols))
		for _, i := range vs.Symbols {
			go func(r chan sflib.StockQuote) {
				// Get Quote
				sq, err := sfc.GetQuote(os.Getenv("STOCKFIGHTER_VENUE"), i.Symbol)
				if err != nil || !sq.Ok {
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					} else {
						fmt.Fprintln(os.Stderr, "Quote request error")
					}
				}
				r <- sq
			}(results)
		}
		select {
		case <-time.After(1 * time.Second):
			if len(results) != 0 {
				fmt.Fprintln(os.Stderr, "Still haven't drained results")
			}
		case r := <-results:
			fmt.Printf("%s@%s B:%d A:%d\n", r.Symbol, r.Venue, r.Bid, r.Ask)
		}
	}
}
