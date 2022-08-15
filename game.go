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
const maxDays int = 2

var stockData = StockData{
	"palantir": StockRange{1, 2},
	"tesla":    StockRange{2, 3},
}
var dailyChoices = []string{"Buy", "Sell", "View details", "Call it a day"}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("\nWelcome to the stocks game!")
	fmt.Printf("The goal is to make as much money in %d days.\n", maxDays)
	fmt.Printf("You start with $%d\n", startingMoney)
	getInput("(press enter to continue...)")

	portfolio := genPortfolio()
	money := startingMoney

	for day := 1; day < (maxDays + 1); day++ {
		portfolio, money = playDay(day, portfolio, money)
	}
	fmt.Println("Thanks for playing!")
}

func getKeys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func genPortfolio() Portfolio {
	portfolio := make(Portfolio)
	for stockName, _ := range stockData {
		portfolio[stockName] = 0
	}
	return portfolio
}

func getInput(question string) (string, error) {
	var input string
	fmt.Print(question)
	_, err := fmt.Scanln(&input)
	return input, err
}

func printOptions(options []string) int {
	fmt.Println("\nOptions:")
	for i, option := range options {
		fmt.Printf("%d: %s\n", i+1, option)
	}
	input, err := getInput("\nChoice: ")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	choice, _ := strconv.Atoi(input)
	return choice
}

func selectOption(options []string) int {
	selection := printOptions(options)
	if selection > len(options)-1 && selection < 1 {
		fmt.Printf("Your choice was invalid. Please choose a number between 1 and %d\n", len(options))
		selection = selectOption(options)
	}

	return selection - 1
}

func genPrices() Portfolio {
	prices := make(Portfolio)
	for stockName, stockRange := range stockData {
		prices[stockName] = rand.Intn(stockRange.max-stockRange.min) + stockRange.min
	}
	return prices
}

func getStockNames(data StockData) []string {
	return getKeys(data)
}

func buy(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	stockNames := getStockNames(stockData)
	// we use stockNames[:] so stockNames are treated as a slice
	stockNumber := selectOption(stockNames[:])
	stockName := stockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := money / stockPrice
	// todo: handle if maxStock is zero
	fmt.Printf("You can buy a max of %d shares", maxStock)
	shares, err := getInput("\nHow many shares?")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	amount, _ := strconv.Atoi(shares)

	cost := stockPrice * amount

	newMoney := money - cost
	portfolio[stockName] = portfolio[stockName] + amount

	return portfolio, newMoney
}

func sell(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	stockNames := getStockNames(stockData)
	stockNumber := selectOption(stockNames[:])
	stockName := stockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := portfolio[stockName]
	// todo: handle if maxStock is zero
	fmt.Printf("You can sell a max of %d shares", maxStock)
	shares, err := getInput("\nHow many shares?")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	amount, _ := strconv.Atoi(shares)

	cost := stockPrice * amount

	newMoney := money + cost
	portfolio[stockName] = portfolio[stockName] - amount

	return portfolio, newMoney
}

func printMap(m map[string]int) {
	// todo: sort by key before printing
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

func printDetails(day int, prices Portfolio, portfolio Portfolio, money int) {
	fmt.Println("\n-------")
	fmt.Printf("Day %d\n", day)
	fmt.Printf("Money: %d\n", money)
	fmt.Println("Portfolio:")
	printMap(portfolio)
	fmt.Println("Prices:")
	printMap(prices)
	fmt.Println("-------")
}

func playDay(day int, portfolio Portfolio, money int) (Portfolio, int) {
	looping := true
	prices := genPrices()
	printDetails(day, prices, portfolio, money)

	for looping {
		switch selectOption(dailyChoices[:]) {
		case 0:
			portfolio, money = buy(prices, portfolio, money)
		case 1:
			portfolio, money = sell(prices, portfolio, money)
		case 2:
			printDetails(day, prices, portfolio, money)
		case 3:
			looping = false
		}
	}
	return portfolio, money
}
