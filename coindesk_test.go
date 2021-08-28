package main

import (
	"context"
	"testing"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestCoindesk(t *testing.T) {
	assert := assert.New(t)

	rate := &CoindeskRate{
		Time: struct {
			Updated time.Time `json:"updatedISO"`
		}{time.Now()},
		Name:       "Dummy",
		Disclaimer: "Dummy",
		BPI: map[string]struct {
			Code   string  `json:"code"`
			Symbol string  `json:"symbol"`
			Rate   float64 `json:"rate_float"`
		}{
			"USD": {
				Code:   "USD",
				Symbol: "$",
				Rate:   48764.234,
			},
		},
	}

	t.Run("API", func(t *testing.T) {
		var r *CoindeskRate

		err := retry.Do(func() error {
			api := NewCoindesk()
			resp, err := api.GetRate(context.Background())
			if err == nil {
				r = resp
			}
			return err
		}, retry.Attempts(10), retry.Delay(1*time.Second))

		if err != nil {
			t.Errorf("error getting coindesk rate: %s\n", err)
		}

		if _, ok := r.BPI["USD"]; !ok {
			t.Errorf("USD is missed in responce")
		}
		if _, ok := r.BPI["EUR"]; !ok {
			t.Errorf("EUR is missed in responce")
		}
		if _, ok := r.BPI["GBP"]; !ok {
			t.Errorf("GBP is missed in responce")
		}
	})

	t.Run("FormatRate", func(t *testing.T) {
		s := rate.FormatRate("USD")
		assert.Equal("$48,764", s)

		s = rate.FormatRate("NotExisting")
		assert.Empty(s)
	})

	t.Run("FormatDesc", func(t *testing.T) {
		s := rate.FormatDesc()
		assert.NotEmpty(s)
	})

}
