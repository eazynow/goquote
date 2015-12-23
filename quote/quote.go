// Package quote contains all business logic aroud creating a load quote.
package quote

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/eazynow/goquote/lender"
)

// Store these as consts for this exercise, however for more flexibility
// these could be passed in also as configuration for a quote
const (
	// MinAmount is the minimum amount allowed for a quote. It is checked by the
	// quote.validate function.
	MinAmount = 1000
	// MaxAmount is the maximum amount allowed for a quote. It is checked by the
	// quote.validate function.
	MaxAmount = 15000
	// compoundFrequency defines the number of periods that interest is compounded
	// to in a year
	compoundFrequency = 12.0
)

// quote represents the structure for a quote response. It is private and needs
// to be creeated using the NewQuote function.
type quote struct {
	lenders          lender.Lenders
	loanPeriodMonths int
	RequestedAmount  int
	Rate             float64
	MonthlyRepayment float64
	TotalRepayment   float64
}

// validate is a private function that validates the quote request criteria for
// preset rules.
// It returns an error detailing the validation failure if one is found. If
// there are no errors then nil is returned.
func (q *quote) validate() error {

	// Check the quote rate is greater than the minimum amount and reject if not.
	if q.RequestedAmount < MinAmount {
		return errors.New(
			fmt.Sprintf(
				"Loan amount of £%d is too low. Minimum loan amount is £%d",
				q.RequestedAmount,
				MinAmount))
	}

	// Check the quote rate is less than the maximum amount and reject if not.
	if q.RequestedAmount > MaxAmount {
		return errors.New(
			fmt.Sprintf(
				"Loan amount of £%d is too high. Maximum loan amount is £%d",
				q.RequestedAmount,
				MaxAmount))
	}

	// Check that the amount is a multiple of 100 and reject if not.
	if q.RequestedAmount%100 != 0 {
		return errors.New("Loan amount must be a multiple of £100")
	}

	return nil
}

// calculate is a private function that processes the quote request and
// calculates the monthly rate and APR.
// Results are populated into the quote structure.
// Returns any error found when attempting to calculate the quote. If there are
// no errors then nil is returned.
func (q *quote) calculate() error {
	// sort lenders into order of preference ascending order (based on rate)
	sort.Sort(q.lenders)

	balance := q.RequestedAmount

	// convert the load period into a float ready for rate calculation.
	// You cannot mix integers and floats in computation.
	fPeriod := float64(q.loanPeriodMonths)

	blendedRate := 0.0

	// Loop through the sorted lender list.
	for _, l := range q.lenders {

		// Find out how much we can borrow from this lender.
		amount := l.Borrow(balance)
		fAmount := float64(amount)

		// Calculate the monthly repayment for the lender.
		q.MonthlyRepayment += calculateMonthlyRate(l.Rate, fPeriod, fAmount)

		blendedRate += fAmount * l.Rate

		balance -= amount

		if balance == 0 {
			break
		}
	}

	// Outstanding balance means there was insufficient funds available from the pool
	// of lenders.
	if balance > 0 {
		return errors.New("It is not possible to provide a quote at this time.")
	}

	// Calculate final values for the quote
	q.TotalRepayment = fPeriod * q.MonthlyRepayment
	q.Rate = blendedRate / float64(q.RequestedAmount)

	// No errors so quote has been calculated!
	return nil
}

// String is a public function to return a string represenatation of a quote results.
// Used to satisfy the fmt.Stringer interface.
// Returns a string representing the quote results.
func (q *quote) String() string {

	// Build a slice up containing the return format.
	var s []string
	s = append(s, fmt.Sprintf("Requested amount: £%d", q.RequestedAmount))
	s = append(s, fmt.Sprintf("Rate: %.1f%%", q.Rate*100.0))
	s = append(s, fmt.Sprintf("Monthly repayment: £%.2f", q.MonthlyRepayment))
	s = append(s, fmt.Sprintf("Total repayment: £%.2f", q.TotalRepayment))

	// Return the string, separated by line breaks.
	return strings.Join(s, "\n")
}

// NewQuote is the main function used to generate a quote based on the amount
// provided, the loanPeriod in months and a slice of lenders.
// The request is quote request is first validated against business rules
// and then calculated.
// Returns a pointer to a quote structure with the quote details if successful.
// If unsuccesful due to validation or calculation issues then an error is returned.
func NewQuote(amount, loanPeriod int, lenders lender.Lenders) (*quote, error) {

	quote := quote{
		RequestedAmount:  amount,
		lenders:          lenders,
		loanPeriodMonths: loanPeriod,
	}

	// Attempt to validate the quote input variables. If validation fails an erroe
	// is returned and passed back to calling function.
	if err := quote.validate(); err != nil {
		return nil, err
	}

	// Attempt to calculate the quote input variables. If validation fails an erroe
	// is returned and passed back to calling function.
	if err := quote.calculate(); err != nil {
		return nil, err
	}

	return &quote, nil
}

// calculateMonthlyRate contains the formula for calculating a compound interest monthly
// rate based on the annualRate, loanPeriod in months and amount required.
// The formula is based on the PMT() excel function:
// monthly_amount = amount * (month_rate*(1+month_rate)^num_periods)/((1+month_rate)^num_periods-1)
// Returns the monthly rate.
func calculateMonthlyRate(annualRate, loanPeriod, amount float64) float64 {

	compoundRate := annualRate / compoundFrequency

	// Apply the formula and return the monthly rate.
	return amount * (compoundRate *
		math.Pow((1.0+compoundRate), loanPeriod)) / (math.Pow(1.0+compoundRate, loanPeriod) - 1)
}
