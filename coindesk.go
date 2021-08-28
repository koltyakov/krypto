package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

// Coindesk API struct
type Coindesk struct{}

// CoindeskRate API response struct
type CoindeskRate struct {
	Time struct {
		Updated time.Time `json:"updatedISO"`
	} `json:"time"`
	Name       string `json:"chartName"`
	Disclaimer string `json:"disclaimer"`
	BPI        map[string]struct {
		Code   string  `json:"code"`
		Symbol string  `json:"symbol"`
		Rate   float64 `json:"rate_float"`
	} `json:"bpi"`
}

// NewCoindesk constructor
func NewCoindesk() Coindesk {
	return Coindesk{}
}

// GetRate gets rate from the API
func (c Coindesk) GetRate(ctx context.Context) (*CoindeskRate, error) {
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

	price := &CoindeskRate{}
	if err := json.Unmarshal(bodyBytes, &price); err != nil {
		return nil, err
	}

	return price, nil
}

// FormatRate formats provided currency rate
func (r *CoindeskRate) FormatRate(currency string) string {
	c, ok := r.BPI[currency]
	if !ok {
		return ""
	}

	return html.UnescapeString(c.Symbol) + formatNumber(int(c.Rate), ',')
}

// FormatDesc formats description (tooltip details message)
func (r *CoindeskRate) FormatDesc() string {
	formatLine := func(curr string) string {
		if _, ok := r.BPI[curr]; !ok {
			return ""
		}
		return "- " + curr + ": " + r.FormatRate(curr) + "\n"
	}
	return ("Updated at " + r.Time.Updated.Local().Format(time.RFC1123) + ":\n" +
		formatLine("USD") + formatLine("EUR") + formatLine("GBP") +
		"\nDisclaimer: " + r.Disclaimer)
}
