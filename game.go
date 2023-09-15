package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Portfolio = map[string]int
type StockRange struct {
	min int
	max int
}
type StockData = map[string]StockRange
type HighScoreEntry struct {
	score int
	date  string
}

const startingMoney int = 111
const maxDays int = 7
const leaderboardFile = "leaderboard.txt"
const maxHighScores = 5

var stockData = StockData{
	"palantir":  StockRange{7, 20},
	"ford":      StockRange{8, 14},
	"microsoft": StockRange{180, 280},
}
var dailyChoices = []string{"Buy", "Sell", "View details", "Call it a day"}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("\nWelcome to the stocks game!")
	fmt.Printf("The goal is to make as much money in %d days.\n", maxDays)
	fmt.Printf("You start with $%d\n", startingMoney)
	fmt.Println("\n(HINT) If you ever want to buy/sell the max amount of a stock, enter 'm' as the amount of the stock you want to buy/sell")
	getInput("\n(press enter to continue...)")

	portfolio := genPortfolio()
	money := startingMoney

	for day := 1; day < (maxDays + 1); day++ {
		portfolio, money = playDay(day, portfolio, money)
		if day == maxDays-1 {
			fmt.Println("\nThis is the last day... sell everything you have!")
		}
	}
	fmt.Println("\nThanks for playing!")
	fmt.Printf("You ended with $%d ($%d total profit)\n", money, money-startingMoney)
	recordScore(money)
}

func recordScore(money int) {
	highScores := readHighScores()
	newHighScores := updateHighScores(highScores, money)
	writeHighScores(newHighScores)
}

func readHighScores() []HighScoreEntry {
	file, err := os.Open(leaderboardFile)
	var lines []HighScoreEntry
	if err != nil {
		fmt.Println("Error opening file:", err)
		return lines
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		intScore, _ := strconv.Atoi(line[0])
		v := HighScoreEntry{intScore, line[1]}
		lines = append(lines, v)
	}
	return lines
}

func insertSlice[T any](slice []T, element T, index int) []T {
	// First, check if the index is out of bounds.
	if index < 0 || index > len(slice) {
		fmt.Println("Index out of bounds.")
		return slice
	}

	// Create a new slice with an extra element.
	newSlice := make([]T, len(slice))

	// Copy the elements before the index.
	copy(newSlice[:index], slice[:index])

	// Insert the element at the specified index.
	newSlice[index] = element

	// Copy the elements after the index.
	copy(newSlice[index+1:], slice[index:maxHighScores-1])

	return newSlice
}

func updateHighScores(highScores []HighScoreEntry, score int) []HighScoreEntry {
	for i, v := range highScores {
		highScore := v.score
		if score > highScore {
			currentTime := time.Now()
			newHighScore := HighScoreEntry{score, fmt.Sprintf("(%s)", currentTime.Format("2006-01-02"))}
			return insertSlice(highScores, newHighScore, i)
		}
	}
	return highScores
}

func writeHighScores(highScores []HighScoreEntry) {
	file, err := os.Create("leaderboard.txt")
	if err != nil {
		fmt.Println("Error creating/opening file:", err)
		return
	}
	defer file.Close() // Close the file when we're done

	for _, score := range highScores {
		line := fmt.Sprintf("%d %s\n", score.score, score.date)
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
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

func printOptions(options []string, sorted bool) int {
	if sorted {
		sort.Strings(options)
	}

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

func selectOption(options []string, sorted bool) int {
	selection := printOptions(options, sorted)
	if selection > len(options)-1 && selection < 1 {
		fmt.Printf("Your choice was invalid. Please choose a number between 1 and %d\n", len(options))
		selection = selectOption(options, sorted)
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

func getBuyableStockNames(stockNames []string, prices Portfolio, money int) []string {
	var buyableStockNames []string
	for _, name := range stockNames {
		if prices[name] <= money {
			buyableStockNames = append(buyableStockNames, name)
		}
	}
	return buyableStockNames
}

func buy(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	stockNames := getStockNames(stockData)
	buyableStockNames := getBuyableStockNames(stockNames, prices, money)
	stockNumber := selectOption(buyableStockNames[:], true)
	stockName := buyableStockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := money / stockPrice
	fmt.Printf("You can buy a max of %d shares", maxStock)
	shares, err := getInput("\nHow many shares? ")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	amount := 0
	if shares == "m" {
		amount = maxStock
	} else {
		amount, _ = strconv.Atoi(shares)
	}

	cost := stockPrice * amount

	newMoney := money - cost
	portfolio[stockName] = portfolio[stockName] + amount

	return portfolio, newMoney
}

func getSellableStockNames(stockNames []string, portfolio Portfolio) []string {
	var sellableStockNames []string
	for _, name := range stockNames {
		if portfolio[name] > 0 {
			sellableStockNames = append(sellableStockNames, name)
		}
	}
	return sellableStockNames
}

func sell(prices Portfolio, portfolio Portfolio, money int) (Portfolio, int) {
	stockNames := getStockNames(stockData)
	sellableStockNames := getSellableStockNames(stockNames, portfolio)
	stockNumber := selectOption(sellableStockNames[:], true)
	stockName := sellableStockNames[stockNumber]
	stockPrice := prices[stockName]
	maxStock := portfolio[stockName]
	fmt.Printf("You can sell a max of %d shares", maxStock)
	shares, err := getInput("\nHow many shares? ")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	amount := 0
	if shares == "m" {
		amount = maxStock
	} else {
		amount, _ = strconv.Atoi(shares)
	}

	cost := stockPrice * amount

	newMoney := money + cost
	portfolio[stockName] = portfolio[stockName] - amount

	return portfolio, newMoney
}

func printMap(m map[string]int) {
	mapKeys := getKeys(m)
	sort.Strings(mapKeys)
	for _, k := range mapKeys {
		fmt.Printf("  %s: %d\n", k, m[k])
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
		switch selectOption(dailyChoices[:], false) {
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
