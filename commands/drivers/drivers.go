package drivers

import (
	"github.com/statecrafthq/borg/utils"
)

// IDFunction is a function for ID resolvers
type IDFunction = func(feature *utils.Feature) ([]string, error)

// ExtrasFunction builds extra fields
type ExtrasFunction func(feature *utils.Feature, extras *Extras) error

// RecordFunction provides record type for each feature
type RecordFunction func(feature *utils.Feature) (RecordType, error)

type RetiredFunction func(feature *utils.Feature) (RetiredType, error)

// EmptyExtras is a function for empty extras
func EmptyExtras(feature *utils.Feature, extras *Extras) error {
	return nil
}

// DefaultRecord is a function for providing default record type (Primary)
func DefaultRecord(feature *utils.Feature) (RecordType, error) {
	return Primary, nil
}

// IgnoreWithoutGeometry is a default behaviour that ignores all records without geometry
func IgnoreWithoutGeometry(feature *utils.Feature) (RecordType, error) {
	if feature.Geometry == nil {
		return Ignored, nil
	}
	return Primary, nil
}

// NoRetired is a default behaviour for retired flag
func NoRetired(feature *utils.Feature) (RetiredType, error) {
	return Active, nil
}

type RecordType int32

type RetiredType int32

const (
	Auxiliary RecordType = 0
	Primary   RecordType = 1
	Ignored   RecordType = 2
)

const (
	Unkwnon RetiredType = 0
	Retired RetiredType = 1
	Active  RetiredType = 2
)

// Driver handles all configuration for specific driver
type Driver struct {
	ID      IDFunction
	Extras  ExtrasFunction
	Record  RecordFunction
	Retired RetiredFunction
}

// Drivers provides all available drivers
func Drivers() map[string]Driver {
	res := make(map[string]Driver)
	res["nyc_blocks"] = NewYorkBlocksDriver()
	res["nyc_parcels"] = NewYorkParcelsDriver()
	res["sf_lots"] = SanFranciscoLotsDriver()
	return res
}
