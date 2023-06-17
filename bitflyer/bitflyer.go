package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const baseURL = "https://api.bitflyer.jp/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	return &APIClient{
		key:        key,
		secret:     secret,
		httpClient: &http.Client{},
	}
}
func (c *APIClient) header(method, endpoint, body string) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)

	sign := hmac.New(sha256.New, []byte(c.secret))
	sign.Write([]byte(message))
	signStr := hex.EncodeToString(sign.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       c.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      signStr,
		"Content-Type":     "application/json",
	}
}

func (c *APIClient) do(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	endpoint, err := url.JoinPath(baseURL, urlPath)
	if err != nil {
		return nil, err
	}
	log.Printf("action=do endpoint=%v", endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range c.header(method, req.URL.RequestURI(), string(data)) {
		req.Header.Add(k, v)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

func (c *APIClient) Balance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := c.do("GET", url, nil, nil)

	log.Printf("url=%v resp=%v", url, string(resp))
	if err != nil {
		log.Printf("action=Balance err=%v", err.Error())
		return nil, err
	}

	var balance []Balance
	if err = json.Unmarshal(resp, &balance); err != nil {
		log.Printf("action=Balance err=%v", err.Error())
		return nil, err
	}
	return balance, err
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (t *Ticker) MidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=DateTime err=%v", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

func (c *APIClient) Ticker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := c.do("GET", url, map[string]string{"product_code": productCode}, nil)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err = json.Unmarshal(resp, &ticker); err != nil {
		return nil, err
	}
	return &ticker, err
}

type JSONRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	ID      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

func (c *APIClient) ReacTimeTicker(symbol string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	if err := conn.WriteJSON(&JSONRPC2{
		Version: "2.0",
		Method:  "subscribe",
		Params:  &SubscribeParams{Channel: "lightning_ticker_" + symbol},
	}); err != nil {
		log.Fatal("subscribe:", err)
	}

LOOP:
	for {
		var res JSONRPC2
		if err := conn.ReadJSON(&res); err != nil {
			log.Println("read:", err)
			return
		}
		if res.Method == "channelMessage" {
			switch v := res.Params.(type) {
			case map[string]interface{}:
				for k, bin := range v {
					if k == "message" {
						marshaTic, err := json.Marshal(bin)
						if err != nil {
							continue LOOP
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue LOOP
						}
						ch <- ticker
					}
				}
			}
		}
	}
}
