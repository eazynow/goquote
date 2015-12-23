package lender

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBorrowWorksWhenEnoughFunds(t *testing.T) {
	l := Lender{}

	l.Available = 600

	amountToBorrow := 500

	amountCanBorrow := l.Borrow(amountToBorrow)

	assert.Equal(
		t,
		amountToBorrow,
		amountCanBorrow,
		"Expected to be able to borrow the full amount")
}

func TestBorrowWorksWhenNotEnoughFunds(t *testing.T) {
	l := Lender{}

	l.Available = 300

	amountToBorrow := 500

	amountCanBorrow := l.Borrow(amountToBorrow)

	assert.Equal(
		t,
		l.Available,
		amountCanBorrow,
		"Expected to be able to borrow 300 only")
}

func TestLendersSortWorksOnRate(t *testing.T) {
	var lenders Lenders

	low := Lender{"low", 0.03, 100}
	mid := Lender{"mid", 0.05, 100}
	high := Lender{"high", 0.08, 100}

	lenders = append(lenders, mid)
	lenders = append(lenders, high)
	lenders = append(lenders, low)

	// this will test Len, Less and Swap functions in one go
	sort.Sort(lenders)

	assert.Equal(t, low, lenders[0], "low rate lender should be first")
	assert.Equal(t, mid, lenders[1], "mid rate lender should be second")
	assert.Equal(t, high, lenders[2], "high rate lender should be third")

}

func TestLendersSortWorksOnAvailable(t *testing.T) {
	var lenders Lenders

	low := Lender{"low", 0.03, 100}
	mid := Lender{"mid", 0.05, 200}
	high := Lender{"high", 0.05, 100}

	lenders = append(lenders, mid)
	lenders = append(lenders, high)
	lenders = append(lenders, low)

	// this will test Len, Less and Swap functions in one go
	sort.Sort(lenders)

	assert.Equal(t, low, lenders[0], "low rate lender should be first")
	assert.Equal(t, mid, lenders[1], "mid lender should be second as has more funds available")
	assert.Equal(t, high, lenders[2], "high lender should be third as he has smallest pool")

}
