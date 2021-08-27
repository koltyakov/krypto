package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"github.com/koltyakov/krypto/icon"
)

var (
	appname = "krypto"
	version string
)

var (
	appConf              = &settings{}
	menu                 = map[string]*systray.MenuItem{}
	tray                 = &Tray{} // Tray state cache
	appCtx, appCtxCancel = context.WithCancel(context.Background())
	currentPrice         *CurrentPrice
	currentCurrency      = "USD" // ToDo: Get from settings
	errorsCount          = 0
)

// Init systray applications
func main() {
	// Graceful shutdown signalling
	grace := make(chan os.Signal, 1)
	signal.Notify(grace, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Lock session to prevent multiple simultaneous application instances
	if err := lockSession(); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer unlockSession()

	// Systray exit handler
	onExit := func() {
		appCtxCancel()
	}

	// Graceful shutdown action
	go func() {
		<-grace
		systray.Quit()
	}()

	// Initiate systray application
	systray.Run(onReady, onExit)
}

// onReady bootstraps system tray menu logic
func onReady() {
	tray.SetIcon(icon.Base)
	tray.SetTitle(" Loading...")

	// Get app settings
	c, err := getSettings()
	if err != nil {
		onError(err)
	}
	appConf = &c

	// Menu items
	menu["currency"] = systray.AddMenuItem("Currency", "")
	menu["currency:usd"] = menu["currency"].AddSubMenuItem("USD", "United States Dollar")
	menu["currency:eur"] = menu["currency"].AddSubMenuItem("EUR", "Euro")
	menu["currency:gbp"] = menu["currency"].AddSubMenuItem("GBP", "British pound sterling")
	systray.AddSeparator()
	menu["about"] = systray.AddMenuItem("About", "Krypto v"+getAppVersion())
	menu["quit"] = systray.AddMenuItem("Quit", "Quit Krypto")

	changeCurrency(currentCurrency) // ToDo: Use settings

	// Menu actions
	go menuActions()

	// Infinite service loop
	for {
		<-time.After(run(1*time.Second, appConf))
	}
}

// menuActions watch to menu actions channels
// must be started in a goroutine, otherwise blocks the loop
func menuActions() {
	for {
		select {
		case <-menu["currency:usd"].ClickedCh:
			onCurrencyChange("usd")
		case <-menu["currency:eur"].ClickedCh:
			onCurrencyChange("eur")
		case <-menu["currency:gbp"].ClickedCh:
			onCurrencyChange("gbp")
		// case <-menu["settings"].ClickedCh:
		// 	openSettingsHandler()
		case <-menu["about"].ClickedCh:
			onOpenLinkHandler("https://github.com/koltyakov/krypto")
		case <-menu["quit"].ClickedCh:
			systray.Quit()
			return
		}
	}
}

// run executes notification checks logic
func run(timeout time.Duration, cnfg *settings) time.Duration {
	price, err := getCurretPrice(appCtx)
	if err != nil {
		// Coindesk API sometimes fails, reducing UI error appearence cases with retries
		errorsCount++
		// On initial load service error retry every 5 seconds but not more than 10 times
		if currentPrice == nil && errorsCount <= 10 {
			return 5 * time.Second
		}
		// Show error only after 10 retries
		if errorsCount > 10 {
			onError(err)
		}
		// On error, awaiting a bit more before quering API again
		return 60 * time.Second
	}

	errorsCount = 0
	currentPrice = price
	onPriceChange()
	return timeout
}

// onCurrencyChange notification mode (all, favorite) change handler
func onCurrencyChange(currency string) {
	changeCurrency(currency)
}

// onError system tray menu on error event handler
func onError(err error) {
	fmt.Printf("error: %s\n", err)
	tray.SetTitle(" Error")
	tray.SetTooltip(fmt.Sprintf("Error: %s", err))
	// tray.SetIcon(icon.Err)
}

// onOpenLinkHandler on open link handler
func onOpenLinkHandler(url string) {
	if err := openBrowser(url); err != nil {
		fmt.Printf("error opening browser: %s\n", err)
	}
}

// changeCurrency sets check mark a currency item
// unchecks other selected modes
func changeCurrency(currency string) {
	currentCurrency = strings.ToUpper(currency)
	for mKey, mItem := range menu {
		if strings.Contains(mKey, "currency:") {
			if mKey == "currency:"+strings.ToLower(currency) {
				mItem.Check()
				continue
			}
			mItem.Uncheck()
		}
	}
	onPriceChange()
}

// onPriceChange updates price in tray
func onPriceChange() {
	if currentPrice == nil {
		return
	}

	rate := currentPrice.FormatRate(currentCurrency)
	if rate == "" {
		onError(fmt.Errorf("can't resolve currency: %s", currentCurrency))
		return
	}

	tray.SetTitle(" " + rate)
	tray.SetTooltip(currentPrice.FormatDescription())
}
