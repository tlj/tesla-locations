package console

import (
	"fmt"
	"tesla-locations/app"
	"tesla-locations/locations"
)

type Console struct {
	formatStr string
	app       *app.App
}

func NewOutput(a *app.App) *Console {
	return &Console{
		app:       a,
	}
}

func (c *Console) Open() error {
	fmt.Printf("%-40s %-15s %-10s %-15s %-15s %-5s %10s\n", "Title", "City", "Country", "Latitude", "Longitude", "Cost", "Unit")

	return nil
}

func (c *Console) Store(location locations.Location) error {
	costPerUnit := 0.0
	billingType := app.BillingTypePerKwh

	if cost, ok := c.app.Prices[location.Country]; ok {
		costPerUnit = cost.ConvertedCostPerUnit
		billingType = cost.BillingType
	}

	fmt.Printf("%-40s %-15s %-10s %-15s %-15s %-5f %10s\n",
		location.Title, location.City, location.Country, location.Latitude, location.Latitude, costPerUnit, billingType)

	return nil
}

func (c *Console) Close() error {
	return nil
}

func (c *Console) Name() string {
	return "Console"
}
