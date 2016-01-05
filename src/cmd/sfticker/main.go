package main

import (
	"fmt"
	"os"
	"sflib"
	"time"
)

func main() {
	sfc := sflib.NewStockfighterClient(os.Getenv("STARFIGHTER_KEY"))
	hb, err := sfc.GetHeartbeat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stderr, "Heartbeat API is OK", hb.Ok)
	if !hb.Ok {
		fmt.Fprintln(os.Stderr, hb.Error)
		os.Exit(1)
		return
	}
	fmt.Fprintln(os.Stderr, "Raw venue is", os.Getenv("STOCKFIGHTER_VENUE"))
	cv, err := sfc.CheckVenue(os.Getenv("STOCKFIGHTER_VENUE"))
	if err != nil || !cv.Ok {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connectivity broken", err)
		}
		fmt.Fprintln(os.Stderr, "Venue isn't up")
	}
	fmt.Fprintln(os.Stderr, "Venue is ", cv.Venue)
	// Get stocks
	vs, err := sfc.GetVenueStocks(os.Getenv("STOCKFIGHTER_VENUE"))
	if err != nil || !vs.Ok {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connectivity broken", err)
		}
		fmt.Fprintln(os.Stderr, "Venue isn't up")
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
