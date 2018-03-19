package drivers

import (
	"github.com/statecrafthq/borg/utils"
)

// IDFunction is a function for ID resolvers
type IDFunction = func(feature *utils.Feature) ([]string, error)

// ExtrasFunction Builds extra fields
type ExtrasFunction func(feature *utils.Feature, extras *Extras) error

type RecordFunction func(feature *utils.Feature) (RecordType, error)

// EmptyExtras is a function for empty extras
func EmptyExtras(feature *utils.Feature, extras *Extras) error {
	return nil
}

func DefaultRecord(feature *utils.Feature) (RecordType, error) {
	return Primary, nil
}

func IgnoreWithoutGeometry(feature *utils.Feature) (RecordType, error) {
	if feature.Geometry == nil {
		return Ignored, nil
	}
	return Primary, nil
}

type RecordType int32

const (
	Auxiliary RecordType = 0
	Primary   RecordType = 1
	Ignored   RecordType = 2
)

// Driver handles all configuration for specific driver
type Driver struct {
	ID     IDFunction
	Extras ExtrasFunction
	Record RecordFunction
}

// Drivers provides all available drivers
func Drivers() map[string]Driver {
	res := make(map[string]Driver)
	res["nyc_blocks"] = NewYorkBlocksDriver()
	res["nyc_parcels"] = NewYorkParcelsDriver()
	res["sf_lots"] = SanFranciscoLotsDriver()
	return res
}
