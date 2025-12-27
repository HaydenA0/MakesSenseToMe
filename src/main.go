package main

import (
	"flag"
	"fmt"
	"strconv"
)

var MINUTE = 60
var HOUR = 3600
var DAY = 86400
var WEEK = 604800
var MONTH = 2592000
var YEAR = 31536000
var CENTURY = 3153600000

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

func handleScalesTime(scales []int) ([]int, error) {

	var scalesTime []int
	if scales[0] == 0 {
		return []int{}, fmt.Errorf("Cant compare to a zero (Division by Zero)")
	}
	for _, scale := range scales {
		inTime := scale / scales[0]
		scalesTime = append(scalesTime, inTime)
	}
	return scalesTime, nil
}

func timeScaleDecider(scalesTime []int) string {
	biggestTime := scalesTime[len(scalesTime)-1]
	if biggestTime <= 100*MINUTE {
		return "Minutes"
	} else if biggestTime <= 100*HOUR {
		return "Hours"
	} else if biggestTime <= 100*DAY {
		return "Days"
	} else if biggestTime <= 100*WEEK {
		return "Weeks"
	} else if biggestTime <= 100*MONTH {
		return "Months"
	} else if biggestTime <= 100*YEAR {
		return "Years"
	} else if biggestTime <= 100*CENTURY {
		return "Centuries"
	} else {
		return "IDK lol"
	}
}

func printScaleComparaison(scalesTime []int, timeScale string) {
	var denominatorFactor int
	switch timeScale {
	case "Minutes":
		denominatorFactor = MINUTE
	case "Hours":
		denominatorFactor = HOUR
	case "Days":
		denominatorFactor = DAY
	case "Weeks":
		denominatorFactor = WEEK
	case "Months":
		denominatorFactor = MONTH
	case "Years":
		denominatorFactor = YEAR
	case "Centuries", "Centruries":
		denominatorFactor = CENTURY
	default:
		denominatorFactor = 1
	}

	fmt.Printf("\nRelative Comparison Table (Unit: %s)\n", timeScale)
	fmt.Println("-------------------------------------------")
	for i, seconds := range scalesTime {
		scaledValue := float64(seconds) / float64(denominatorFactor)
		if scaledValue < 0.01 {
			fmt.Printf("Scale %d: %10.4f %s\n", i+1, scaledValue, timeScale)
		} else {
			fmt.Printf("Scale %d: %10.2f %s\n", i+1, scaledValue, timeScale)
		}
	}
	fmt.Println("-------------------------------------------")
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
		scalesTime, err := handleScalesTime(scales[:])
		if err != nil {
			fmt.Printf("Error in handleScalesTime : ")
			fmt.Println(err)
			return
		}
		fmt.Println("Debug : ", scalesTime)
		timeScale := timeScaleDecider(scalesTime[:])
		printScaleComparaison(scalesTime, timeScale)

	case "distance":
		fmt.Printf("You chose %s ?\n", *pScaleType)
	}
}
