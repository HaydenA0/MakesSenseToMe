package main

import (
	"flag"
	"fmt"
	"strconv"
)

type Unit struct {
	name        string
	valueInBase float64
	limit       int
}

var timeUnits = []Unit{
	{"second", 1, 60},
	{"minute", 60, 60},
	{"hour", 3600, 24},
	{"day", 86400, 7},
	{"week", 604800, 4},
	{"month", 2592000, 12},
	{"year", 31536000, 100},
	{"decade", 315360000, 10},
	{"century", 3153600000, 10},
	{"millennium", 31536000000, 10},
	{"geological_era", 3153600000000, 10},
	{"age_of_earth", 1.38e17, 10},
	{"age_of_universe", 4.35e17, 10},
}

var distanceUnits = []Unit{
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

func inputToFloat(prompt string) float64 {
	var buffer string
	fmt.Printf("%s", prompt)
	fmt.Scanln(&buffer)
	o, err := strconv.ParseFloat(buffer, 64)
	if err != nil {
		fmt.Print("Error converting, enter an integer\n")
	}
	return o

}

func prettyPrint(scales_index float64, scaleSeconds float64, unit Unit) {
	fmt.Println("+------------------+----------------------------+")
	fmt.Printf("| %-16s | %-26v |\n", "OriginalValue", scales_index)
	fmt.Printf("| %-16s | %-26.2f |\n", "ScaleValue", scaleSeconds/float64(unit.valueInBase))
	fmt.Printf("| %-16s | %-26s |\n", "UnitName", unit.name)
	fmt.Println("+------------------+----------------------------+")
}
func normalizeScale(scales []float64) []float64 {
	normalizedScale := make([]float64, len(scales))
	for index := range scales {
		normalizedScale[index] = float64(scales[index]) / float64(scales[0])
	}
	return normalizedScale
}

func getBestUnit(unitSlice []Unit, scaleSeconds float64) Unit {
	for _, unit := range unitSlice {
		if scaleSeconds <= float64(unit.limit)*float64(unit.valueInBase) {
			return unit
		}
	}
	return unitSlice[len(unitSlice)-1]
}

func handleScales(scales []float64, unitSlice []Unit) {
	normalizedScale := normalizeScale(scales)
	for index, scaleSeconds := range normalizedScale {
		unit := getBestUnit(unitSlice, scaleSeconds)
		prettyPrint(scales[index], scaleSeconds, unit)
	}
}

func main() {
	pnumArgs := flag.Int("na", 2, "Number of scales to compare")
	pScaleType := flag.String("s", "time", "Type of scale [time, distance]")
	flag.Parse()
	var scales []float64
	for i := 0; i < *pnumArgs; i++ {
		// TODO : make this in a function and easy to use
		prompt := fmt.Sprintf("Enter the %dth scale (increasing order) : > ", i)
		scale := inputToFloat(prompt)
		scales = append(scales, scale)
	}
	switch *pScaleType {
	case "time":
		handleScales(scales[:], timeUnits[:])

	case "distance":
		handleScales(scales[:], distanceUnits[:])
	}
}
