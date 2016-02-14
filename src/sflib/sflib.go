package sflib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Heartbeat struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type CheckVenue struct {
	Ok    bool   `json:"ok"`
	Venue string `json:"venue"`
	Error string `json:"error"`
}

type StockSymbol struct {
	Name, Symbol string
}

type VenueStocks struct {
	Ok      bool `json:"ok"`
	Venue   *string
	Symbols []StockSymbol
	Error   string `json:"error"`
}

type Quote struct {
	Symbol, Venue                                                  string
	Bid, Ask, BidSize, AskSize, BidDepth, AskDepth, Last, LastSize int
	LastTrade, QuoteTime                                           time.Time
}

// https://starfighter.readme.io/docs/a-quote-for-a-stock
type StockQuote struct {
	Ok bool
	Quote
}

// https://starfighter.readme.io/docs/quotes-ticker-tape-websocket
// every-so-slightly different interface than the polling quote interface,
// for whatever reason- nested quote object instead of embedded
type StreamingStockQuote struct {
	Ok    bool
	quote Quote
}

type OrderType string

const (
	Market     = OrderType("market")
	Limit      = OrderType("limit")
	FillOrKill = OrderType("fill-or-kill")
	Immediate  = OrderType("immediate-or-cancel")
)

type Direction string

const (
	Buy  = Direction("buy")
	Sell = Direction("sell")
)

type Order struct {
	Account   string    `json:"account"`
	Venue     string    `json:"venue"`
	Stock     string    `json:"stock"`
	Direction Direction `json:"direction"`
	OrderType OrderType `json:"ordertype"`
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
}

type Fill struct {
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
	Ts    string `json:"ts"`
}

type OrderResponse struct {
	Ok          bool      `json:"ok"`
	Open        bool      `json:"open"`
	Error       string    `json:"error"`
	Account     string    `json:"account"`
	Venue       string    `json:"venue"`
	Ts          string    `json:"ts"` // ISO-8601 Timestamp. Not ready to commit to actually parsing this yet ...
	Symbol      string    `json"symbol"`
	Id          int       `json:"id"`
	TotalFilled int       `json:"totalFilled"`
	Price       int       `json:"price"`
	OriginalQty int       `json:"originalQty"`
	Qty         int       `json:"qty"`
	Direction   Direction `json:"direction"`
	OrderType   OrderType `json:"orderType"`
	Fills       []Fill    `json:"fills"`
}

type StockfighterClient struct {
	httpclient *http.Client
	Api_key    string
}

func NewStockfighterClient(api_key string) *StockfighterClient {
	sfc := &StockfighterClient{
		Api_key: api_key,
	}
	sfc.httpclient = &http.Client{}
	return sfc
}

func newHttpReq(method, url string, body *interface{}) (*interface{}, error) {
	return nil, nil
}

func (this *StockfighterClient) GetHeartbeat() (Heartbeat, error) {
	var hb Heartbeat
	req, err := http.NewRequest("GET", "https://api.stockfighter.io/ob/api/heartbeat", nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		return hb, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hb, err
	}
	err = json.Unmarshal(body, &hb)
	if err != nil {
		return hb, err
	}
	return hb, err
}

func (this *StockfighterClient) CheckVenue(venue string) (CheckVenue, error) {
	var cv CheckVenue
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/heartbeat", venue),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		return cv, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cv, err
	}
	err = json.Unmarshal(body, &cv)
	if err != nil {
		return cv, err
	}
	return cv, err
}

func (this *StockfighterClient) GetVenueStocks(venue string) (VenueStocks, error) {
	var vs VenueStocks
	vs.Venue = &venue
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks", venue),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		return vs, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vs, err
	}
	err = json.Unmarshal(body, &vs)
	if err != nil {
		return vs, err
	}
	return vs, err
}

type OrderStatus struct {
	Ok          bool      `json:"ok"`
	Error       *string   `json:"error"`
	Symbol      string    `json:"symbol"`
	Venue       string    `json:"venue"`
	Direction   Direction `json:"direction"`
	OriginalQty int       `json:"originalQty"`
	Qty         int       `json:"qty"`
	Price       int       `json:"price"`
	OrderType   OrderType `json:"orderType"`
	Id          int       `json:"id"`
	Account     string    `json:"account"`
	Ts          time.Time `json:"ts"`
	Fills       []Fill    `json:"fills"`
	TotalFilled int       `json:"totalFilled"`
	Open        bool      `json:"open"`
}

type StockOrders struct {
	Ok     bool          `json:"ok"`
	Venue  string        `json:"venue"`
	Error  string        `json:"error"`
	Orders []OrderStatus `json:"orders"`
}

func (this *StockfighterClient) CheckAllOrderStatus(venue, account, stock string) (StockOrders, error) {
	var so StockOrders
	fmt.Println("Venue ", venue, " account ", account, " stock ", stock)
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/accounts/%s/stocks/%s/orders", venue, account, stock),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "HTTP Req failed", err)
		return so, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading the body of the reply failed", err)
		return so, err
	}
	err = json.Unmarshal(body, &so)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshalling failed ", err)
		fmt.Fprintln(os.Stderr, "Got ", string(body))
		return so, err
	}
	if so.Error != "" {
		return so, errors.New(fmt.Sprintf("API Error: ", so.Error))
	}
	return so, err
}

func (this *StockfighterClient) CancelOrder(venue, stock string, id int) (OrderStatus, error) {
	var ost OrderStatus
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s/orders/%d", venue, stock, id),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		fmt.Println("Order http request failed")
		return ost, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read order status response")
		return ost, err
	}
	err = json.Unmarshal(body, &ost)
	if err != nil {
		fmt.Println("could not unmarshal order status", string(body))
		return ost, err
	}
	if ost.Error != nil {
		return ost, errors.New(*ost.Error)
	}
	return ost, err
}

func (this *StockfighterClient) CheckOrderStatus(venue, stock string, id int) (OrderStatus, error) {
	var ost OrderStatus
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s/orders/%d", venue, stock, id),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		fmt.Println("Order http request failed")
		return ost, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read order status response")
		return ost, err
	}
	err = json.Unmarshal(body, &ost)
	if err != nil {
		fmt.Println("could not unmarshal order status", string(body))
		return ost, err
	}
	if ost.Error != nil {
		return ost, errors.New(*ost.Error)
	}
	return ost, err
}

func (this *StockfighterClient) PostOrder(venue, symbol, account string, qty, price int, direction Direction, order_type OrderType) (*OrderResponse, error) {
	var order Order
	var response *OrderResponse = &OrderResponse{}
	order.Account = account
	order.Venue = venue
	order.Stock = symbol
	order.Direction = direction
	order.OrderType = order_type
	order.Price = price
	order.Qty = qty

	//fmt.Println("Quantity ", qty, " Price ", price)

	rb, err := json.Marshal(order)
	if err != nil {
		fmt.Sprintf("Could not marshall the order type (!)", err)
		return nil, err
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s/orders", venue, symbol),
		bytes.NewBuffer(rb))
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "http post failed")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "body read failed")
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unmarshalling failed")
		return response, err
	}
	if response.Error != "" {
		fmt.Fprintln(os.Stderr, "got an error from the API", response.Error)
		return response, errors.New(response.Error)
	}

	return response, nil
}

func (this *StockfighterClient) GetQuote(venue, symbol string) (StockQuote, error) {
	var sq StockQuote
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s/quote", venue, symbol),
		nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := this.httpclient.Do(req)
	if err != nil {
		return sq, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sq, err
	}
	err = json.Unmarshal(body, &sq)
	if err != nil {
		return sq, err
	}
	return sq, nil
}

func Setup() (sfc *StockfighterClient, vs *VenueStocks, err error) {
	sfc = NewStockfighterClient(os.Getenv("STARFIGHTER_KEY"))
	hb, err := sfc.GetHeartbeat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Hearbeat failure ", err)
		return
	}
	fmt.Fprintln(os.Stderr, "Heartbeat API is OK", hb.Ok)
	if !hb.Ok {
		fmt.Fprintln(os.Stderr, hb.Error)
		return
	}
	fmt.Fprintln(os.Stderr, "Raw venue is", os.Getenv("STOCKFIGHTER_VENUE"))
	cv, err := sfc.CheckVenue(os.Getenv("STOCKFIGHTER_VENUE"))
	if err != nil || !cv.Ok {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connectivity broken", err)
		} else {
			err = errors.New(cv.Error)
		}
		fmt.Fprintln(os.Stderr, "Venue isn't up")
		return
	}
	fmt.Fprintln(os.Stderr, "Venue is ", cv.Venue)
	// Get stocks
	vs_lit, err := sfc.GetVenueStocks(os.Getenv("STOCKFIGHTER_VENUE"))
	vs = &vs_lit
	if err != nil || !vs.Ok {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connectivity broken", err)
		} else {
			err = errors.New(vs.Error)
		}
		fmt.Fprintln(os.Stderr, "Venue isn't up")
	}
	return
}
