package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/eazynow/goquote/lender"
	"github.com/eazynow/goquote/quote"
)

const (
	// loanPeriodMonths is the default loan length in months.
	loanPeriodMonths = 36
)

var (
	filename string
	amount   int
)

func init() {
	// Skip the first arg as that is the program name.
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Printf("Usage: goquote [filename] [amount]")
		os.Exit(0)
	}

	filename = args[0]

	// Check that the amount provided is a valid integer.
	var err error
	amount, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("The amount %s is not a valid integer", args[1])
		os.Exit(0)
	}
}

func main() {

	// Attempt to import the csv file into a lender.Lenders slice
	lenders, err := lender.ImportCSV(filename)
	if err != nil {
		// There was a problem importing the csv - bail!
		fmt.Println(err)
		os.Exit(0)
	}

	// Attempt to create a new quote based on the input parameters
	q, err := quote.NewQuote(amount, loanPeriodMonths, lenders)

	if err != nil {
		// There was a problem creating the quote - bail!
		fmt.Println(err)
		os.Exit(0)
	}

	// Display the quote
	fmt.Println(q.String())
}
