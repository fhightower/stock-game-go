package main

import (
	"fmt"
)

const startingMoney int = 111
const maxDays int = 1

var stockOptions = [2]string{"palantir", "tesla"}

func main() {
	fmt.Println("Welcome to the stocks game!")
	fmt.Printf("The goal is to make as much money in %d days.", maxDays)
	fmt.Printf("You start with $%d\n", startingMoney)

	// portfolio := make(map[string]int)
	// money := startingMoney

	for day := 1; day < (maxDays + 1); day++ {
		fmt.Printf("Day %d\n", day)
		// portfolio, money = playDay(day, portfolio, money)
	}
}
