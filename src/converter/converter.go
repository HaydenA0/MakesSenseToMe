package converter

type Unit struct {
	Name        string
	ValueInBase float64
	Limit       int
}

type Result struct {
	OriginalValue float64
	ScaledValue   float64
	UnitName      string
}

var TimeUnits = []Unit{
	{"second", 1, 60},
	{"minute", 60, 60},
	{"hour", 3600, 24},
	{"day", 86400, 7},
	{"week", 604800, 4},
	{"month", 2592000, 12},
	{"year", 31536000, 10},
	{"decade", 315360000, 10},
	{"century", 3153600000, 10},
	{"millennium", 31536000000, 10},
	{"myriad", 315360000000, 10},            // 10,000 years
	{"ice_age_cycle", 3153600000000, 10},    // 100,000 years
	{"million_years", 31536000000000, 10},   // 1,000,000 years
	{"billion_years", 31536000000000000, 4}, // 1,000,000,000 years
	{"age_of_earth", 1.38e17, 3},
	{"age_of_universe", 4.35e17, 10},
}

var DistanceUnits = []Unit{
	{"millimeter", 1, 10},
	{"centimeter", 10, 100},
	{"meter", 1000, 1000},
	{"football_field", 100000, 100},
	{"kilometer", 1000000, 1000},
	{"ten_kilometer_city", 10000000, 10},
	{"hundred_kilometer_country", 100000000, 10},
	{"half_the_globe", 10800000000, 10},
	{"spain_to_argentina", 11000000000, 10},
	{"london_to_sydney", 17000000000, 10},
	{"earth_to_moon", 384400000000, 100},
	{"earth_to_sun", 149600000000000, 100},
	{"earth_to_mars", 225000000000000, 100},
	{"width_solar_system", 5900000000000000, 10},
}

func Convert(scales []float64, units []Unit) []Result {
	if len(scales) == 0 {
		return nil
	}

	baseValue := scales[0]
	results := make([]Result, len(scales))

	for i, originalValue := range scales {
		var normalized float64
		if baseValue != 0 {
			normalized = originalValue / baseValue
		}

		bestUnit := getBestUnit(units, normalized)

		var scaledValue float64
		if bestUnit.ValueInBase != 0 {
			scaledValue = normalized / bestUnit.ValueInBase
		}

		results[i] = Result{
			OriginalValue: originalValue,
			ScaledValue:   scaledValue,
			UnitName:      bestUnit.Name,
		}
	}

	return results
}

// ConvertAbsolute treats every input value as an independent absolute quantity
// in the base unit (e.g. seconds, millimetres) and reports each in its best-fit
// unit. Unlike Convert, the first value is not used as a normalisation base.
func ConvertAbsolute(scales []float64, units []Unit) []Result {
	if len(scales) == 0 {
		return nil
	}

	results := make([]Result, len(scales))
	for i, originalValue := range scales {
		bestUnit := getBestUnit(units, originalValue)

		var scaledValue float64
		if bestUnit.ValueInBase != 0 {
			scaledValue = originalValue / bestUnit.ValueInBase
		}

		results[i] = Result{
			OriginalValue: originalValue,
			ScaledValue:   scaledValue,
			UnitName:      bestUnit.Name,
		}
	}
	return results
}

func getBestUnit(unitSlice []Unit, scaleSeconds float64) Unit {
	for _, unit := range unitSlice {
		if scaleSeconds <= float64(unit.Limit)*unit.ValueInBase {
			return unit
		}
	}
	if len(unitSlice) > 0 {
		return unitSlice[len(unitSlice)-1]
	}
	return Unit{}
}
