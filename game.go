package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Portfolio = map[string]int
type StockRange struct {
	min int
	max int
}
type StockData = map[string]StockRange

const startingMoney int = 111
const maxDays int = 1

var stockData = StockData{
	"palantir": StockRange{1, 2},
	"tesla":    StockRange{2, 3},
}
var dailyChoices = []string{"Buy", "Sell", "View details", "Call it a day"}

func main() {
	rand.Seed(time.Now().UnixNano())
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
	fmt.Print(question)
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	return input
}

func printOptions(options []string) int {
	fmt.Println("Options:")
	for i, option := range options {
		fmt.Printf("%d: %s\n", i+1, option)
	}
	choice, _ := strconv.Atoi(getInput("Choice: "))
	return choice
}

func selectOption(options []string) int {
	selection := printOptions(options)
	if selection > len(options)-1 && selection < 1 {
		fmt.Printf("Your choice was invalid. Please choose a number between 1 and %d\n", len(options))
		selection = selectOption(options)
	}

	// todo: once this func is working, update other funcs to use it
	return selection - 1
}

func genPrices() Portfolio {
	prices := make(Portfolio)
	for stockName, stockRange := range stockData {
		prices[stockName] = rand.Intn(stockRange.max-stockRange.min) + stockRange.min
	}
	return prices
}

func getStockNames(data StockData) [2]string {
	var names [2]string
	var i int
	// todo: is it possible to yield in go?
	for stockName, _ := range data {
		names[i] = stockName
		i++
	}
	return names
}

func buy(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	// todo: display the stock names as numbers rather than having to type in the full name
	stockNames := getStockNames(stockData)
	// we use stockNames[:] so stockNames are treated as a slice
	stockNumber := selectOption(stockNames[:])
	stockName := stockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := money / stockPrice
	// todo: handle if maxStock is zero (or negative)
	fmt.Printf("You can buy a max of %d shares", maxStock)
	amount, _ := strconv.Atoi(getInput("How many shares?"))

	cost := stockPrice * amount

	newMoney := money - cost
	portfolio[stockName] = portfolio[stockName] + amount

	return portfolio, newMoney
}

func sell(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	stockNames := getStockNames(stockData)
	// we use stockNames[:] so stockNames are treated as a slice
	stockNumber := selectOption(stockNames[:])
	stockName := stockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := portfolio[stockName]
	// todo: handle if maxStock is zero (or negative)
	fmt.Printf("You can sell a max of %d shares", maxStock)
	amount, _ := strconv.Atoi(getInput("How many shares?"))

	cost := stockPrice * amount

	newMoney := money + cost
	portfolio[stockName] = portfolio[stockName] - amount

	return portfolio, newMoney
}

func printDetails(day int, prices Portfolio, portfolio Portfolio, money int) {
	fmt.Println("-----")
	fmt.Printf("\nDay %d\n", day)
	fmt.Printf("Money: %d\n", money)
	fmt.Printf("Portfolio: %v\n", portfolio)
	fmt.Printf("Prices: %v\n", prices)
}

func playDay(day int, portfolio Portfolio, money int) (Portfolio, int) {
	looping := true
	prices := genPrices()
	printDetails(day, prices, portfolio, money)

	for looping {
		choice := printOptions(dailyChoices[:])
		fmt.Println(choice)
		if choice == 1 {
			portfolio, money = buy(prices, portfolio, money)
		} else if choice == 2 {
			portfolio, money = sell(prices, portfolio, money)
		} else if choice == 3 {
			printDetails(day, prices, portfolio, money)
		} else if choice == 4 {
			looping = false
		}
	}
	return portfolio, money
}
