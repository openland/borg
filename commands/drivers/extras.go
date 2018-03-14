package drivers

func NewExtras() Extras {
	return Extras{Enums: []ExtrasEnum{}, Strings: []ExtrasString{}, Floats: []ExtrasFloat{}, Ints: []ExtrasInt{}}
}

func (e *Extras) HasKey(key string) bool {
	for i := range e.Strings {
		if e.Strings[i].Key == key {
			return true
		}
	}
	for i := range e.Enums {
		if e.Enums[i].Key == key {
			return true
		}
	}
	for i := range e.Floats {
		if e.Floats[i].Key == key {
			return true
		}
	}
	for i := range e.Ints {
		if e.Ints[i].Key == key {
			return true
		}
	}
	return false
}

func (e *Extras) DeleteKey(key string) {

	// Strings
	strings := []ExtrasString{}
	changed := false
	for i := range e.Strings {
		if e.Strings[i].Key != key {
			strings = append(strings, e.Strings[i])
		} else {
			changed = true
		}
	}
	if changed {
		e.Strings = strings
	}

	// Enums
	enums := []ExtrasEnum{}
	changed = false
	for i := range e.Enums {
		if e.Enums[i].Key != key {
			enums = append(enums, e.Enums[i])
		} else {
			changed = true
		}
	}
	if changed {
		e.Enums = enums
	}

	// Floats
	floats := []ExtrasFloat{}
	changed = false
	for i := range e.Floats {
		if e.Floats[i].Key != key {
			floats = append(floats, e.Floats[i])
		} else {
			changed = true
		}
	}
	if changed {
		e.Floats = floats
	}

	// Ints
	ints := []ExtrasInt{}
	changed = false
	for i := range e.Ints {
		if e.Ints[i].Key != key {
			ints = append(ints, e.Ints[i])
		} else {
			changed = true
		}
	}
	if changed {
		e.Ints = ints
	}
}

func (e *Extras) AppendString(key string, value string) {
	e.DeleteKey(key)
	e.Strings = append(e.Strings, ExtrasString{Key: key, Value: value})
}

type Extras struct {
	Enums   []ExtrasEnum   `json:"enums"`
	Strings []ExtrasString `json:"strings"`
	Floats  []ExtrasFloat  `json:"floats"`
	Ints    []ExtrasInt    `json:"ints"`
}

type ExtrasEnum struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}
type ExtrasString struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ExtrasFloat struct {
	Key   string  `json:"key"`
	Value float32 `json:"value"`
}
type ExtrasInt struct {
	Key   string `json:"key"`
	Value int32  `json:"value"`
}
