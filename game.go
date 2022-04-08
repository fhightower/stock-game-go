package main

import (
	"fmt"
	"os"
	"strconv"
)

const startingMoney int = 111
const maxDays int = 1

type Portfolio = map[string]int
type StockRange struct {
	min int
	max int
}
type StockData = map[string]StockRange

var stockData = StockData{
	"palantir": StockRange{1, 2},
	"tesla":    StockRange{2, 3},
}

func main() {
	fmt.Println("Welcome to the stocks game!")
	fmt.Printf("The goal is to make as much money in %d days.\n", maxDays)
	fmt.Printf("You start with $%d\n", startingMoney)

	portfolio := genPortfolio()
	money := startingMoney

	for day := 1; day < (maxDays + 1); day++ {
		portfolio, money = playDay(day, portfolio, money)
	}
}

func genPortfolio() Portfolio {
	portfolio := make(Portfolio)
	for stockName, _ := range stockData {
		portfolio[stockName] = 0
	}
	return portfolio
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
	fmt.Println("1. Buy")
	fmt.Println("2. Sell")
	fmt.Println("3. Show details")
	fmt.Println("4. Call it a day")
	choice, _ := strconv.Atoi(getInput("Choice: "))
	return choice
}

func genPrices() Portfolio {
	// todo: update this
	prices := make(Portfolio)
	return prices
}

func buy(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	// todo: show max number of a stock which can be bought
	// options
	// fmt.Printf("You can buy: %a", options)
	stock := getInput("Which stock?")
	amount, _ := strconv.Atoi(getInput("How many shares?"))

	cost := prices[stock] * amount

	newMoney := money - cost
	portfolio[stock] = portfolio[stock] + amount

	return portfolio, newMoney
}

func playDay(day int, portfolio Portfolio, money int) (Portfolio, int) {
	looping := true
	prices := genPrices()
	fmt.Printf("\nDay %d\n", day)
	fmt.Printf("	Portfolio: %v\n", portfolio)
	fmt.Printf("	Prices: %v\n", prices)

	for looping {
		choice := printOptions()
		fmt.Println(choice)
		if choice == 1 {
			portfolio, money = buy(prices, portfolio, money)
		}
		looping = false
	}
	return portfolio, money
}
