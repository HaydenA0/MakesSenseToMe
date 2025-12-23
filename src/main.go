package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const MAX_GEN_NUMBER = 100

func generateRandomInteger() int {
	number := rand.Intn(MAX_GEN_NUMBER)
	return number
}
func myReadStdio(prompt string, buffer *string) {
	fmt.Print(prompt)
	_, err := fmt.Scanln(buffer)
	if err != nil {
		fmt.Print("Error readin input")
	}
}

func handleUserInput(userGuessInt int, randomNumber int) bool {
	if userGuessInt == randomNumber {
		fmt.Print("Congratulations ! You won\n")
		return true
	} else if userGuessInt > randomNumber {
		fmt.Print("Your guess is greater than the secret number\n")
		return false
	} else {
		fmt.Print("Your guess is smaller than the secret number\n")
		return false
	}
}

func main() {
	randomNumber := generateRandomInteger()
	for true {
		var userGuessString string
		myReadStdio("Input your guess : > ", &userGuessString)
		userGuessInt, err := strconv.Atoi(userGuessString)
		if err != nil {
			print("Error Conversion")
		}
		if handleUserInput(userGuessInt, randomNumber) {
			break
		}
	}

}
