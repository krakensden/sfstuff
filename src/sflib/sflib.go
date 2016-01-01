package sflib

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
)

type Heartbeat struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type StockfighterClient struct {
	httpclient http.Client
	Api_key    string
}

func (this *StockfighterClient) GetHeartbeat() (Heartbeat, error) {
	var hb Heartbeat
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.stockfighter.io/ob/api/heartbeat", nil)
	req.Header.Add("X-Starfighter-Authorization", this.Api_key)
	resp, err := client.Do(req)
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
