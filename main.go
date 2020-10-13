package main

import (
	log "github.com/sirupsen/logrus"
	"tesla-locations/app"
	"tesla-locations/outputs"
	"tesla-locations/outputs/console"
	"tesla-locations/outputs/teslamate_database"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func run(a *app.App) error {
	var output outputs.Output
	switch a.Config.Output {
	case "TeslamateDatabase":
		output = teslamate_database.NewOutput(a)
	case "Console":
		output = console.NewOutput(a)
	default:
		log.Fatalf("Invalid out %s. Valid types are: TeslamateDatabase, Console", a.Config.Output)
	}
	log.Printf("Using output %s.", output.Name())

	if err := output.Open(); err != nil {
		return err
	}
	defer output.Close()

	log.Infof("Updating for countries (%v) and types (%v)...", a.Config.Countries, a.Config.LocationTypes)
	for _, country := range a.Config.Countries {
		price, ok := a.Prices[country]
		if !ok {
			log.Warnf("Warning: %s does not have price information in Config/Prices.json.", country)
		}
		log.Infof("Found price for %s: %f %s %s", country, price.ConvertedCostPerUnit, a.Config.Currency, price.BillingType)
	}

	chargers, err := a.Client.Countries(a.Config.Countries, a.Config.LocationTypes)
	if err != nil {
		return err
	}

	storedCount := 0
	for _, charger := range chargers {
		err = output.Store(charger)
		if err != nil {
			log.Error(err)
		}
		storedCount++
	}
	log.Infof("Update done. %d locations stored.", storedCount)

	return nil
}

func main() {
	a := &app.App{}
	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	log.Infof("Starting TeslaMate Location Geofence Sync...")

	if err := run(a); err != nil {
		log.Fatal(err)
	}

	if a.Config.Daemon && a.Config.IntervalMinutes > 0 {
		log.Infof("Starting daemon, updating every %d minutes.", a.Config.IntervalMinutes)
		ticker := time.NewTicker(time.Duration(a.Config.IntervalMinutes) * time.Minute)
		for {
			select {
			case <-ticker.C:
				if err := run(a); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
