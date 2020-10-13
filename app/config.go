package app

import (
	"encoding/json"
	gf "flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"strings"
	"tesla-locations/locations"
)

type CountryList []string
type LocationTypeList []locations.LocationType

type Config struct {
	Daemon          bool             `json:"daemon" envconfig:"LOCATION_DAEMON"`
	IntervalMinutes int              `json:"interval_minutes" envconfig:"LOCATION_INTERVAL_MINUTES"`
	LogLevel        string           `json:"log_level" envconfig:"LOCATION_LOG_LEVEL"`
	Output          string           `json:"output" envconfig:"LOCATION_OUTPUT"`
	Radius          int              `json:"radius" envconfig:"LOCATION_RADIUS"`
	Currency        string           `json:"currency" envconfig:"LOCATION_CURRENCY"`
	Countries       CountryList      `json:"countries" envconfig:"LOCATION_COUNTRIES"`
	LocationTypes   LocationTypeList `json:"location_types" envconfig:"LOCATION_TYPES"`
}

type BillingType string

const (
	BillingTypePerKwh    BillingType = "per_kwh"
	BillingTypePerMinute BillingType = "per_minute"
)

func (cl *CountryList) Decode(value string) error {
	countries := strings.Split(value, ";")
	out := CountryList{}
	for _, c := range countries {
		out = append(out, c)
	}
	*cl = out
	return nil
}

func (cl *CountryList) String() string {
	var out []string
	for _, c := range *cl {
		out = append(out, c)
	}
	return strings.Join(out, ";")
}

func (cl *CountryList) Set(value string) error {
	*cl = append(*cl, value)
	return nil
}

func (cl *LocationTypeList) Decode(value string) error {
	countries := strings.Split(value, ";")
	out := LocationTypeList{}
	for _, c := range countries {
		out = append(out, locations.LocationType(c))
	}
	*cl = out
	return nil
}

func (cl *LocationTypeList) String() string {
	var out []string
	for _, c := range *cl {
		out = append(out, string(c))
	}
	return strings.Join(out, ";")
}

func (cl *LocationTypeList) Set(value string) error {
	*cl = append(*cl, locations.LocationType(value))
	return nil
}

type Prices map[string]*struct {
	CostPerUnit          float64     `json:"cost_per_unit"`
	Currency             string      `json:"currency"`
	BillingType          BillingType `json:"billing_type"`
	SessionFee           float64     `json:"session_fee"`
	ForeignCurrency      bool
	ConvertedCostPerUnit float64
}

type ExchangeRates struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"Rates"`
}

func (a *App) LoadConfig() error {
	var (
		countries CountryList
		types LocationTypeList
	)
	gf.BoolVar(&a.Config.Daemon, "d", false, "Daemon")
	gf.IntVar(&a.Config.IntervalMinutes, "i", 1440, "Interval (minutes), when in daemon mode")
	gf.IntVar(&a.Config.Radius, "r", 50, "Radius")
	gf.StringVar(&a.Config.LogLevel, "l", "Info", "Log level (debug, info, warn, error)")
	gf.StringVar(&a.Config.Currency, "x", "EUR", "Exchange (Currency)")
	gf.StringVar(&a.Config.Output, "o", "Console", "Output (Console, TeslamateDatabase)")
	gf.Var(&types, "t", "Location type (supercharger, standard charger, destination charger)")
	gf.Var(&countries, "c", "Countries")

	configJson, err := ioutil.ReadFile("config/config.json")
	if err != nil{
		return fmt.Errorf("unable to read file config/config.json: %s", err)
	}

	err = json.Unmarshal(configJson, &a.Config)
	if err != nil {
		return fmt.Errorf("unable to parse file config/config.json: %s", err)
	}

	err = envconfig.Process("", &a.Config)
	if err != nil {
		return err
	}

	gf.Parse()

	if len(countries) > 0 {
		a.Config.Countries = countries
	}

	if len(types) > 0 {
		a.Config.LocationTypes = types
	}

	return nil
}

func (a *App) LoadPrices() error {
	pricesJson, err := ioutil.ReadFile("config/prices.json")
	if err != nil {
		return fmt.Errorf("unable to read file config/prices.json: %s", err)
	}

	err = json.Unmarshal(pricesJson, &a.Prices)
	if err != nil {
		return fmt.Errorf("unable to parse file config/prices.json: %s", err)
	}

	for _, v := range a.Prices {
		v.ConvertedCostPerUnit = v.CostPerUnit
		if v.Currency != a.Config.Currency {
			v.ForeignCurrency = true
			v.ConvertedCostPerUnit, _ = a.ConvertPrice(v.CostPerUnit, v.Currency)
		}
	}

	return nil
}

func (a *App) LoadExchangeRates() error {
	ratesJson, err := ioutil.ReadFile("config/exchange_rates.json")
	if err != nil {
		return fmt.Errorf("unable to read file config/exchange_rates.json: %s", err)
	}

	var exchangeRates ExchangeRates
	err = json.Unmarshal(ratesJson, &exchangeRates)
	if err != nil {
		return fmt.Errorf("unable to parse file config/exchange_rates.json: %s", err)
	}

	a.Rates = exchangeRates.Rates

	return nil
}
