package main

import (
	"flag"
	"fmt"
	"strconv"
)

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
	// turn scales to seconds, raw seconds
	var scalesTime []int
	if scales[0] == 0 {
		return []int{}, fmt.Errorf("Division by zero")
	}
	for _, scale := range scales {
		inTime := scale / scales[0]
		scalesTime = append(scalesTime, inTime)
	}
	return scalesTime, nil
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
		scales_time, err := handleScalesTime(scales[:])
		if err != nil {
			fmt.Printf("Error in handleScalesTime")
			return
		}
		fmt.Println(scales_time)

	case "distance":
		fmt.Printf("You chose %s ?\n", *pScaleType)
	}
}
