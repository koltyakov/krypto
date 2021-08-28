package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// Coingecko API struct
type Coingecko struct{}

// CoingeckoRate API response struct
type CoingeckoRate struct {
	Bitcoin struct {
		USD        int     `json:"usd"`
		USD24hPerc float64 `json:"usd_24h_change"`
		EUR        int     `json:"eur"`
		EUR24hPerc float64 `json:"eur_24h_change"`
		GBP        int     `json:"gbp"`
		GBP24hPerc float64 `json:"gbp_24h_change"`
		Updated    int64   `json:"last_updated_at"`
	} `json:"bitcoin"`
}

// NewCoingecko constructor
func NewCoingecko() Coingecko {
	return Coingecko{}
}

// GetRate gets rate from the API
func (c Coingecko) GetRate(ctx context.Context) (*CoingeckoRate, error) {
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd,eur,gbp&ids=bitcoin&include_24hr_change=true&include_last_updated_at=true", nil)
	if err != nil {
		return nil, err
	}

	ct := http.DefaultTransport.(*http.Transport).Clone()
	// ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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

	rate := &CoingeckoRate{}
	if err := json.Unmarshal(bodyBytes, &rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// FormatRate formats provided currency rate
func (r *CoingeckoRate) FormatRate(currency string) string {
	if currency == "USD" {
		return fmt.Sprintf("$%s %s%.2f%%", formatNumber(r.Bitcoin.USD, ','), trendSymbol(r.Bitcoin.USD24hPerc), math.Abs(r.Bitcoin.USD24hPerc))
	}
	if currency == "EUR" {
		return fmt.Sprintf("€%s %s%.2f%%", formatNumber(r.Bitcoin.EUR, ','), trendSymbol(r.Bitcoin.EUR24hPerc), math.Abs(r.Bitcoin.EUR24hPerc))
	}
	if currency == "GBP" {
		return fmt.Sprintf("₤%s %s%.2f%%", formatNumber(r.Bitcoin.GBP, ','), trendSymbol(r.Bitcoin.GBP24hPerc), math.Abs(r.Bitcoin.GBP24hPerc))
	}
	return ""
}

// FormatDesc formats description (tooltip details message)
func (r *CoingeckoRate) FormatDesc() string {
	formatLine := func(curr string) string {
		if currSrt := r.FormatRate(curr); currSrt != "" {
			return "- " + curr + ": " + currSrt + "\n"
		}
		return ""
	}
	return ("Updated at " + time.Unix(r.Bitcoin.Updated, 0).Local().Format(time.RFC1123) + ":\n" +
		formatLine("USD") + formatLine("EUR") + formatLine("GBP") +
		"\nData source: CoinGecko \nhttps://www.coingecko.com")
}
