package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CurrentPrice struct {
	Time struct {
		UpdatedISO time.Time `json:"updatedISO"`
	} `json:"time"`
	ChartName string                  `json:"chartName"`
	BPI       map[string]CurrencyRate `json:"bpi"`
}

type CurrencyRate struct {
	Code   string  `json:"code"`
	Symbol string  `json:"symbol"`
	Rate   float32 `json:"rate_float"`
}

func getCurretPrice(ctx context.Context) (*CurrentPrice, error) {
	req, err := http.NewRequest("GET", "https://api.coindesk.com/v1/bpi/currentprice.json", nil)
	if err != nil {
		return nil, err
	}

	ct := http.DefaultTransport.(*http.Transport).Clone()
	ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Transport: ct, Timeout: 10 * time.Second}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (Status Code %d)", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	price := &CurrentPrice{}
	if err := json.Unmarshal(bodyBytes, &price); err != nil {
		return nil, err
	}

	return price, nil
}
