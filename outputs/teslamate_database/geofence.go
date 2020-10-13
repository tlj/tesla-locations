package teslamate_database

import (
	"context"
	"database/sql"
	"fmt"
	"tesla-locations/app"
	"time"

	_ "github.com/lib/pq"
)

type Geofence struct {
	Id          int
	SourceId    string
	Name        string
	Latitude    float64
	Longitude   float64
	Radius      int
	InsertedAt  time.Time
	UpdatedAt   time.Time
	CostPerUnit float64
	SessionFee  float64
	BillingType app.BillingType // per_kwh, per_minute
}

type RepositoryInterface interface {
	Insert(ctx context.Context, g *Geofence) error
	Upsert(ctx context.Context, g *Geofence) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Upsert(ctx context.Context, g *Geofence) error {
	query := `SELECT id FROM geofences WHERE name LIKE '%(' || $1 || ')' LIMIT 1`
	res, err := r.db.QueryContext(ctx, query, g.SourceId)
	if err != nil {
		return err
	}
	defer res.Close()

	if !res.Next() {
		return r.Insert(ctx, g)
	}

	var id int64
	_ = res.Scan(&id)

	query = `UPDATE geofences SET name = $2, latitude = $3, longitude = $4, updated_at = $5, cost_per_unit = $6, session_fee = $7, billing_type = $8 WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, id, fmt.Sprintf("%s (%s)", g.Name, g.SourceId), g.Latitude, g.Longitude, time.Now(), g.CostPerUnit, g.SessionFee, g.BillingType)

	return nil
}

func (r *repository) Insert(ctx context.Context, g *Geofence) error {
	query := `INSERT INTO geofences 
				(name, latitude, longitude, radius, inserted_at, updated_at, cost_per_unit, session_fee, billing_type) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, fmt.Sprintf("%s (%s)", g.Name, g.SourceId), g.Latitude, g.Longitude, g.Radius, g.InsertedAt, g.UpdatedAt, g.CostPerUnit, g.SessionFee, g.BillingType)
	if err != nil {
		return err
	}

	return nil
}
