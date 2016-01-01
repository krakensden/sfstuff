package sflib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Heartbeat struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type CheckVenue struct {
	Ok    bool   `json:"ok"`
	Venue string `json:"venue"`
}

type StockSymbol struct {
	Name, Symbol string
}

type VenueStocks struct {
	Ok      bool `json:"ok"`
	Symbols []StockSymbol
}

type StockQuote struct {
	Ok                                                             bool
	Symbol, Venue                                                  string
	Bid, Ask, BidSize, AskSize, BidDepth, AskDepth, Last, LastSize int
	LastTrade, QuoteTime                                           time.Time
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
