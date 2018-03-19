package drivers

import "github.com/statecrafthq/borg/utils"

func sanFranciscoLotsID(feature *utils.Feature) ([]string, error) {
	return []string{feature.Properties["blklot"].(string)}, nil
}

func sanFranciscoClassifier(feature *utils.Feature) (RecordType, error) {
	if feature.Properties["blklot"] != feature.Properties["mapblklot"] {
		return Auxiliary, nil
	}
	return Primary, nil
}

func SanFranciscoLotsDriver() Driver {
	return Driver{ID: sanFranciscoLotsID, Extras: EmptyExtras, Record: sanFranciscoClassifier}
}
