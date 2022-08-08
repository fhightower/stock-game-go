package main

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenPrices(t *testing.T) {
	assert := assert.New(t)
	got := genPrices()
	log.Print(got)
	assert.Equal(len(got), len(stockData), "Prices are generated for each stock")
}

