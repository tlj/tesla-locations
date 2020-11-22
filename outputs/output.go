package outputs

import (
	"github.com/tlj/tesla"
)

type Output interface {
	Open() error
	Store(location tesla.Location) error
	Close() error
	Name() string
}
