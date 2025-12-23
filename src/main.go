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

func main() {

	pnumArgs := flag.Int("na", 2, "Number of scales to compare")
	flag.Parse()
	var scales []int
	for i := 0; i < *pnumArgs; i++ {
		prompt := fmt.Sprintf("Enter the %dth scale : > ", i)
		scale := inputInt(prompt)
		scales = append(scales, scale)
	}

}
