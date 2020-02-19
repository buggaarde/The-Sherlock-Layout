package main

import (
	"fmt"

	"github.com/buggaarde/layout-optimizer/layout"
	"github.com/buggaarde/layout-optimizer/text"
)

func main() {
	fs := text.AnalyseDir("./texts")

	keyb := layout.Load("sherlock.toml")
	bestKeyb := keyb
	fmt.Printf("Initial layout is the following:\n\n")
	layout.Print(bestKeyb)
	layout.AnalyseLayout(bestKeyb, fs)
	layout.SingleHandUtilization(bestKeyb, "./texts")

	temps := []float64{
		25000,
		15000,
		10000, 10000, 10000,
		7000, 7000, 7000,
		5000, 5000, 5000,
		3000, 3000, 3000,
		1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000,
	}

	for _, initialTemp := range temps {
		bestKeyb = layout.Anneal(bestKeyb, initialTemp, fs.T, fs.D, fs.S)
		layout.Write(bestKeyb, "currently_best.toml")
	}

	fmt.Printf("The best layout found by this process is:\n\n")
	layout.Print(bestKeyb)
	layout.AnalyseLayout(bestKeyb, fs)
	layout.SingleHandUtilization(bestKeyb, "./texts")
}
