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
	api_key    string
}

func (this *StockfighterClient) GetHeartbeat() {
}

func GetHeartbeat() (Heartbeat, error) {
	//resp, err := http.Get("http://example.com/")
	//resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	//resp, err := http.PostForm("http://example.com/form",
	//		url.Values{"key": {"Value"}, "id": {"123"}})

	var hb Heartbeat
	resp, err := http.Get("https://api.stockfighter.io/ob/api/heartbeat")
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
