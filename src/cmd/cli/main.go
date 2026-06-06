package main

import (
	"flag"
	"fmt"
	"make/sense/to/me/src/converter"
	"strconv"
)

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

func prettyPrint(res converter.Result) {
	fmt.Println("+------------------+----------------------------+")
	fmt.Printf("| %-16s | %-26v |\n", "OriginalValue", res.OriginalValue)
	fmt.Printf("| %-16s | %-26.2f |\n", "ScaleValue", res.ScaledValue)
	fmt.Printf("| %-16s | %-26s |\n", "UnitName", res.UnitName)
	fmt.Println("+------------------+----------------------------+")
}

func handleScaleType(scaleType *string, scales []float64) {
	var results []converter.Result

	switch *scaleType {
	case "time":
		results = converter.Convert(scales, converter.TimeUnits)
	case "distance":
		results = converter.Convert(scales, converter.DistanceUnits)
	}

	for _, res := range results {
		prettyPrint(res)
	}
}

func fillScales(pnumArgs *int) []float64 {
	var scales []float64
	for i := 0; i < *pnumArgs; i++ {
		prompt := fmt.Sprintf("Enter the %dth scale (increasing order) : > ", i+1)
		scale := inputToFloat(prompt)
		scales = append(scales, scale)
	}
	return scales
}

func getFlags() (*int, *string) {
	numScales := flag.Int("na", 2, "Number of scales to compare")
	scaleType := flag.String("s", "time", "Type of scale [time, distance]")
	return numScales, scaleType
}

func main() {
	numScales, scaleType := getFlags()

	flag.Parse()

	scales := fillScales(numScales)

	handleScaleType(scaleType, scales)
}
