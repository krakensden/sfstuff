package main

import (
	"fmt"
	"os"
	"sflib"
	"strconv"
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
	target_price, err := strconv.Atoi(os.Getenv("PRICE"))
	if target == "" || account == "" || err != nil {
		fmt.Println(os.Stderr, "Need STOCK, ACCOUNT, and PRICE environment variables to implement a buying strategy")
		os.Exit(2)
	}

	var orders []int
	//var filled_qty, last_filled_qty int = 0, 0

	no_such_stock := true
	for _, symbol := range stocks.Symbols {
		if target == symbol.Symbol {
			no_such_stock = false
		}
	}

	if no_such_stock {
		fmt.Fprint(os.Stderr, "Stock ", target, " is not on the exchange, just ")
		for _, symbol := range stocks.Symbols {
			fmt.Fprint(os.Stderr, symbol, ", ")
		}
		fmt.Fprintln(os.Stderr)
	}

	err = clear_old_orders(client, *stocks.Venue, account, target)
	if err != nil {
		os.Exit(5)
	}

	for {
		// check
		results := make(chan sflib.StockQuote, len(stocks.Symbols))
		go func(r chan sflib.StockQuote) {
			// Get Quote
			sq, err := client.GetQuote(*stocks.Venue, target)
			if err != nil || !sq.Ok {
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Fprintln(os.Stderr, "Quote request error")
				}
			}
			r <- sq
		}(results)

		select {
		case <-time.After(1 * time.Second):
			if len(results) != 0 {
				fmt.Fprintln(os.Stderr, "Still haven't drained results")
			}
		case r := <-results:

			orders, _ = cull_dead_orders(client, *stocks.Venue, target, orders)

			if len(orders) > 1 {
				continue
			}

			if (r.Ask < target_price && r.Ask != 0) || (r.Last < target_price) {
				order_type := sflib.Limit

				price := target_price
				oresp, err := client.PostOrder(*stocks.Venue, target, account, 50, price, sflib.Buy, order_type)
				orders = append(orders, oresp.Id)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Print(".")
				}
			}
		}
	}
}

func clear_old_orders(client *sflib.StockfighterClient, venue string, account string, stock string) error {
	sq, err := client.CheckAllOrderStatus(venue, account, stock)
	for _, order := range sq.Orders {
		if order.Open {
			fmt.Println("Clearing order ", order)
			client.CancelOrder(venue, stock, order.Id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
		}
	}
	return err
}
func cull_dead_orders(client *sflib.StockfighterClient, venue string, stock string, orders []int) ([]int, int) {
	var rval []int
	qty_purchased := 0

	for _, order_id := range orders {
		if order_id == -1 {
			continue
		}
		// check order
		orderstatus, err := client.CheckOrderStatus(venue, stock, order_id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if !orderstatus.Open {
			for _, f := range orderstatus.Fills {
				qty_purchased += f.Qty
			}
		} else {
			rval = append(rval, order_id)
		}
	}
	return rval, qty_purchased
}
