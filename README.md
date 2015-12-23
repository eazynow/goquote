# goquote - Compound Interest Rate Calculator in Go


[![Build Status](https://travis-ci.org/eazynow/goquote.svg)](https://travis-ci.org/eazynow/goquote)

Go code (golang) commmand line utility for calculating monthly compound interest rates in go based on multiple lenders.

## Requirements

goquote should work with go versions 1.2 or greater.

## Installation

To install goquote, use `go get`:
```
go get github.com/eazynow/goquote
```

## Usage

To use the utility, run goquote from the command line passing in the csv file containing the lender pool and the amount to borrow.

```
$ ./goquote market.csv 1000
Requested amount: £1000
Rate: 7.0%
Monthly repayment: £30.88
Total repayment: £1111.64
```

## Tests

To run the tests, use `go test`:
```
go test ./...
```

