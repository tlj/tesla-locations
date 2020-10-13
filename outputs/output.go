package outputs

import (
	"tesla-locations/locations"
)

type Output interface {
	Open() error
	Store(location locations.Location) error
	Close() error
	Name() string
}
