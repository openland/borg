package drivers

import (
	"github.com/statecrafthq/borg/utils"
)

func sanFranciscoLotsID(feature *utils.Feature) ([]string, error) {
	if feature.Properties["blklot"] != feature.Properties["mapblklot"] {
		return []string{feature.Properties["mapblklot"].(string), feature.Properties["blklot"].(string)}, nil
	}
	return []string{feature.Properties["mapblklot"].(string)}, nil
}

func sanFranciscoClassifier(feature *utils.Feature) (RecordType, error) {
	// Just for now ignore everything without geometry
	if feature.Geometry == nil {
		return Ignored, nil
	}
	if feature.Properties["blklot"] != feature.Properties["mapblklot"] {
		return Auxiliary, nil
	}
	return Primary, nil
}

func sanFranciscoRetired(feature *utils.Feature) (RetiredType, error) {
	mdrop, mdropPresent := feature.Properties["mad_drop"]
	rdrop, rdropPresent := feature.Properties["rec_drop"]
	if mdropPresent || rdropPresent {
		if mdrop != nil || rdrop != nil {
			return Retired, nil
		}
		return Active, nil
	}
	return Unkwnon, nil
}

func SanFranciscoLotsDriver() Driver {
	return Driver{ID: sanFranciscoLotsID, Extras: EmptyExtras, Record: sanFranciscoClassifier, Retired: sanFranciscoRetired, MultipleID: false}
}
