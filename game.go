package main

import (
	"fmt"
	"os"
	"strconv"
)

const startingMoney int = 111
const maxDays int = 1

type Portfolio = map[string]int
type StockNames = [2]string

var stockOptions = StockNames{"palantir", "tesla"}

func main() {
	fmt.Println("Welcome to the stocks game!")
	fmt.Printf("The goal is to make as much money in %d days.", maxDays)
	fmt.Printf("You start with $%d\n", startingMoney)

	portfolio := make(Portfolio)
	money := startingMoney

	for day := 1; day < (maxDays + 1); day++ {
		fmt.Printf("Day %d\n", day)
		portfolio, money = playDay(day, portfolio, money)
	}
}

func getInput(question string) string {
	var input string
	fmt.Println(question)
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	return input
}

func printOptions() int {
	fmt.Println("Options:")
	fmt.Println("1. Show details")
	fmt.Println("2. Buy")
	fmt.Println("3. Sell")
	fmt.Println("4. Call it a day")
	choice := getInput("Choice: ")
	i, _ := strconv.Atoi(choice)
	return i
}

func playDay(day int, portfolio Portfolio, money int) (Portfolio, int) {
	looping := true

	for looping {
		choice := printOptions()
		fmt.Println(choice)
		looping = false
	}
	return portfolio, money
}
