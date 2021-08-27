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

type CurrentPrice struct {
	Time struct {
		Updated time.Time `json:"updatedISO"`
	} `json:"time"`
	Name       string                  `json:"chartName"`
	Disclaimer string                  `json:"disclaimer"`
	BPI        map[string]CurrencyRate `json:"bpi"`
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

func (p *CurrentPrice) FormatRate(currency string) string {
	c, ok := p.BPI[currency]
	if !ok {
		return ""
	}

	return html.UnescapeString(c.Symbol) + formatNumber(int(c.Rate), ',')
}

func (p *CurrentPrice) FormatDescription() string {
	return ("Updated at " + p.Time.Updated.Local().Format(time.RFC1123) + ":\n" +
		"- USD: " + html.UnescapeString(p.BPI["USD"].Symbol) + formatNumber(int(p.BPI["USD"].Rate), ',') + "\n" +
		"- EUR: " + html.UnescapeString(p.BPI["EUR"].Symbol) + formatNumber(int(p.BPI["EUR"].Rate), ',') + "\n" +
		"- GBP: " + html.UnescapeString(p.BPI["GBP"].Symbol) + formatNumber(int(p.BPI["GBP"].Rate), ',') + "\n\n" +
		"Disclaimer: " + p.Disclaimer)
}
