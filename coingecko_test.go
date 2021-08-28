package main

import (
	"context"
	"testing"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestCoingecko(t *testing.T) {
	assert := assert.New(t)

	rate := &CoingeckoRate{
		Bitcoin: struct {
			USD        int     `json:"usd"`
			USD24hPerc float64 `json:"usd_24h_change"`
			EUR        int     `json:"eur"`
			EUR24hPerc float64 `json:"eur_24h_change"`
			GBP        int     `json:"gbp"`
			GBP24hPerc float64 `json:"gbp_24h_change"`
			Updated    int64   `json:"last_updated_at"`
		}{
			USD:        48764,
			USD24hPerc: 0.2389640295086352,
			EUR:        41181,
			EUR24hPerc: -0.0389640295086352,
			Updated:    1630164081,
		},
	}

	t.Run("API", func(t *testing.T) {
		var r *CoingeckoRate

		err := retry.Do(func() error {
			api := NewCoingecko()
			resp, err := api.GetRate(context.Background())
			if err == nil {
				r = resp
			}
			return err
		}, retry.Attempts(1), retry.Delay(1*time.Second))

		if err != nil {
			t.Errorf("error getting coindesk rate: %s\n", err)
		}

		assert.NotEmpty(r.Bitcoin.USD)
		assert.NotEmpty(r.Bitcoin.EUR)
		assert.NotEmpty(r.Bitcoin.GBP)
	})

	t.Run("FormatRate", func(t *testing.T) {
		s := rate.FormatRate("USD")
		assert.Equal("$48,764 ↑0.24%", s)

		s = rate.FormatRate("EUR")
		assert.Equal("€41,181 ↓0.04%", s)

		s = rate.FormatRate("NotExisting")
		assert.Empty(s)
	})

	t.Run("FormatDesc", func(t *testing.T) {
		s := rate.FormatDesc()
		assert.NotEmpty(s)
	})

}
