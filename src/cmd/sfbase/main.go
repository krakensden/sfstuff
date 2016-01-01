package main

import (
	"fmt"
	"os"
	"sflib"
)

func main() {
	sfc := sflib.NewStockfighterClient(os.Getenv("STARFIGHTER_KEY"))
	hb, err := sfc.GetHeartbeat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Heartbeat API is OK", hb.Ok)
	if !hb.Ok {
		fmt.Println(hb.Error)
		os.Exit(1)
		return
	}
	cv, err := sfc.CheckVenue(os.Getenv("STOCKFIGHTER_VENUE"))
	if err != nil || !cv.Ok {
		if err != nil {
			fmt.Println("Connectivity broken", err)
		}
		fmt.Println("Venue isn't up")
	}
	fmt.Println("Venue is ", cv.Venue)
	// Get stocks
	vs, err := sfc.GetVenueStocks(os.Getenv("STOCKFIGHTER_VENUE"))
	if err != nil || !vs.Ok {
		if err != nil {
			fmt.Println("Connectivity broken", err)
		}
		fmt.Println("Venue isn't up")
	}
	fmt.Println("Available stocks")
	for _, i := range vs.Symbols {
		fmt.Println("Venue is ", i.Name, " :: ", i.Symbol)
	}
}
