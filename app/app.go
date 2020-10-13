package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"tesla-locations/locations"
)

type App struct {
	Config Config
	Client locations.ClientInterface
	Prices Prices
	Rates  map[string]float64
}

func (a *App) ConvertPrice(from float64, fromCurrency string) (float64, error) {
	if fromCurrency == a.Config.Currency {
		return from, nil
	}

	finalRate, ok := a.Rates[a.Config.Currency]
	if !ok {
		log.Warnf("Configured base currency %s is not available in our exchange Rates.", a.Config.Currency)
		return 0, fmt.Errorf("invalid base currency %s", a.Config.Currency)
	}

	rate, ok := a.Rates[fromCurrency]
	if !ok {
		log.Warnf("The currency %s is not available in our exchange Rates.", fromCurrency)
		return 0, fmt.Errorf("invalid currency %s", fromCurrency)
	}

	// convert to base currency (EUR)
	eur := from / rate

	// convert to requested currency
	out := eur * finalRate

	return out, nil
}

func (a *App) Init() error {
	if err := a.LoadConfig(); err != nil {
		return err
	}

	switch a.Config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	if err := a.LoadExchangeRates(); err != nil {
		return err
	}

	if err := a.LoadPrices(); err != nil {
		return err
	}

	a.Client = locations.NewClient("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0")

	return nil
}
