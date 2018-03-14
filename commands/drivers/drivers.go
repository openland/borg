package drivers

import (
	"github.com/statecrafthq/borg/utils"
)

// IDFunction is a function for ID resolvers
type IDFunction = func(feature *utils.Feature) ([]string, error)

// ExtrasFunction Builds extra fields
type ExtrasFunction func(feature *utils.Feature, extras *Extras) error

// EmptyExtras is a function for empty extras
func EmptyExtras(feature *utils.Feature, extras *Extras) error {
	return nil
}

// Driver handles all configuration for specific driver
type Driver struct {
	ID     IDFunction
	Extras ExtrasFunction
}

// Drivers provides all available drivers
func Drivers() map[string]Driver {
	res := make(map[string]Driver)
	res["nyc_blocks"] = NewYorkBlocksDriver()
	res["nyc_parcels"] = NewYorkParcelsDriver()
	return res
}
