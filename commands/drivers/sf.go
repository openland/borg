package drivers

import (
	"strings"

	"github.com/statecrafthq/borg/utils"
)

func sanFranciscoLotsID(feature *utils.Feature) ([]string, error) {
	if feature.Properties["blklot"] != feature.Properties["mapblklot"] {
		return []string{strings.TrimLeft(feature.Properties["mapblklot"].(string), "0"), strings.TrimLeft(feature.Properties["blklot"].(string), "0")}, nil
	}
	return []string{strings.TrimLeft(feature.Properties["mapblklot"].(string), "0")}, nil
}

func sanFranciscoBlocksID(feature *utils.Feature) ([]string, error) {
	return []string{strings.TrimLeft(feature.Properties["block_num"].(string), "0")}, nil
}

//block_num

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
	return Driver{ID: sanFranciscoLotsID, Extras: EmptyExtras, Record: sanFranciscoClassifier, Retired: sanFranciscoRetired}
}

func SanFranciscoBlocksDriver() Driver {
	return Driver{ID: sanFranciscoBlocksID, Extras: EmptyExtras, Record: DefaultRecord, Retired: NoRetired}
}
