package main

import (
	"fmt"
	"os"
	"sflib"
)

func main() {
	sfc := &sflib.StockfighterClient{
		Api_key: os.Getenv("STARFIGHTER_KEY"),
	}
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
}
