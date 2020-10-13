package teslamate_database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"strconv"
	"tesla-locations/app"
	"tesla-locations/locations"
	"time"
)

type Output struct {
	db   *sql.DB
	repo RepositoryInterface
	app  *app.App
}

func NewOutput(app *app.App) *Output {
	return &Output{
		app: app,
	}
}

type dbConfigEnv struct {
	DatabaseUser string `envconfig:"DATABASE_USER" required:"true"`
	DatabasePass string `envconfig:"DATABASE_PASS" required:"true"`
	DatabaseName string `envconfig:"DATABASE_NAME" required:"true"`
	DatabaseHost string `envconfig:"DATABASE_HOST" required:"true"`
}

func (t *Output) Open() error {
	var dbConfig dbConfigEnv
	err := envconfig.Process("", &dbConfig)
	if err != nil {
		return err
	}

	dbConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbConfig.DatabaseUser,
		dbConfig.DatabasePass,
		dbConfig.DatabaseHost,
		dbConfig.DatabaseName)

	t.db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		return fmt.Errorf("unable to open db connection: %s", err)
	}

	err = t.db.Ping()
	if err != nil {
		return fmt.Errorf("unable to ping db: %s", err)
	}

	t.repo = NewRepository(t.db)

	return nil
}

func (t *Output) Store(location locations.Location) error {
	// Don't save chargers which haven't opened yet
	if location.OpenSoon == "1" {
		return nil
	}
	lat, _ := strconv.ParseFloat(location.Latitude, 64)
	long, _ := strconv.ParseFloat(location.Longitude, 64)
	costPerUnit := 0.0
	billingType := app.BillingTypePerKwh
	sessionFee := 0.0

	if cost, ok := t.app.Prices[location.Country]; ok {
		costPerUnit = cost.ConvertedCostPerUnit
		billingType = cost.BillingType
	}

	g := Geofence{
		Id:          0,
		Name:        fmt.Sprintf("%s, %s, %s", location.Title, location.City, location.Country),
		SourceId:    location.Nid,
		Latitude:    lat,
		Longitude:   long,
		Radius:      t.app.Config.Radius,
		InsertedAt:  time.Now(),
		UpdatedAt:   time.Now(),
		CostPerUnit: costPerUnit,
		SessionFee:  sessionFee,
		BillingType: billingType,
	}
	err := t.repo.Upsert(context.Background(), &g)
	if err != nil {
		log.Errorf("Error while inserting %s: %s.", location.Title, err)
	}
	log.Debugf("Stored %s.", g.Name)

	return nil
}

func (t *Output) Close() error {
	return t.db.Close()
}

func (t *Output) Name() string {
	return "TeslamateDatabase"
}
