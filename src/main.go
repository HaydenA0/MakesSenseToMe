package main

import (
	"flag"
	"fmt"
	"strconv"
)

type Unit struct {
	name           string
	valueInSeconds int
	limit          int
}

var timeUnits = []Unit{
	{"second", 1, 60},
	{"minute", 60, 60},
	{"hour", 3600, 24},
	{"day", 86400, 7},
	{"week", 604800, 4},
	{"month", 2592000, 12},
	{"year", 31536000, 100},
	{"century", 3153600000, 10},
}

func inputInt(prompt string) int {
	var buffer string
	fmt.Printf("%s", prompt)
	fmt.Scanln(&buffer)
	o, err := strconv.Atoi(buffer)
	if err != nil {
		fmt.Print("Error converting, enter an integer\n")
	}
	return o

}

func prettyPrint(scales_index int, scaleSeconds float32, unit Unit) {
	fmt.Println("+------------------+----------------------------+")
	fmt.Printf("| %-16s | %-26v |\n", "OriginalValue", scales_index)
	fmt.Printf("| %-16s | %-26.2f |\n", "ScaleValue", scaleSeconds/float32(unit.valueInSeconds))
	fmt.Printf("| %-16s | %-26s |\n", "UnitName", unit.name)
	fmt.Println("+------------------+----------------------------+")
}
func handleTimeScale(scales []int) {
	normalizedScale := make([]float32, len(scales))
	for index := range scales {
		normalizedScale[index] = float32(scales[index]) / float32(scales[0])
	}
	fmt.Println("DEBUG : normalizedScale : ", normalizedScale)
	for index, scaleSeconds := range normalizedScale {
		for _, unit := range timeUnits {
			if scaleSeconds <= float32(unit.limit)*float32(unit.valueInSeconds) {
				prettyPrint(scales[index], scaleSeconds, unit)
				break
			}
		}
	}
}

func main() {
	pnumArgs := flag.Int("na", 2, "Number of scales to compare")
	pScaleType := flag.String("s", "time", "Type of scale [time, distance]")
	flag.Parse()
	var scales []int
	for i := 0; i < *pnumArgs; i++ {
		prompt := fmt.Sprintf("Enter the %dth scale (increasing order) : > ", i)
		scale := inputInt(prompt)
		scales = append(scales, scale)
	}
	switch *pScaleType {
	case "time":
		handleTimeScale(scales[:])

	case "distance":
		fmt.Printf("You chose %s ?\n", *pScaleType)
	}
}
