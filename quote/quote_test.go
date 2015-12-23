package quote

import (
	"fmt"
	"strings"
	"testing"

	"github.com/eazynow/goquote/lender"
	"github.com/stretchr/testify/assert"
)

func TestQuoteValidateFailsOnLowAmount(t *testing.T) {
	q := quote{}

	// use the constant set so if it changes, the test adapts
	q.RequestedAmount = MinAmount - 1

	err := q.validate()

	assert.NotNil(t, err, "Expected an error as requested amount was too low")

}

func TestQuoteValidateFailsOnHighAmount(t *testing.T) {
	q := quote{}

	// use the constant set so if it changes, the test adapts
	q.RequestedAmount = MaxAmount + 1

	err := q.validate()

	assert.NotNil(t, err, "Expected an error as requested amount was too high")
}

func TestQuoteValidateFailsOnMultiple(t *testing.T) {
	q := quote{}

	// use the constant set so if it changes, the test adapts
	q.RequestedAmount = MinAmount + 10

	err := q.validate()

	assert.NotNil(t, err, "Expected an error as requested amount was not divisible by 100")
}

func TestQuoteValidateWorksOnGoodValues(t *testing.T) {
	q := quote{}

	// use the constant set so if it changes, the test adapts
	q.RequestedAmount = MinAmount + 100

	err := q.validate()

	assert.Nil(t, err, "Didn't expect an error as requested amount was divisible by 100 and in range")
}

func TestQuoteCalculateFailsWhenAmountNotAvailable(t *testing.T) {
	amount := 1300

	// only 1200 available in total
	var lenders lender.Lenders
	l1 := lender.Lender{"", 0.01, 1000}
	l2 := lender.Lender{"", 0.01, 200}

	lenders = append(lenders, l1)
	lenders = append(lenders, l2)

	q := quote{}
	q.RequestedAmount = amount
	q.lenders = lenders
	q.loanPeriodMonths = 36

	err := q.calculate()

	assert.NotNil(t, err, "Expected an error as pool amount is insufficient")
}

func TestQuoteCalculateWorksWithCorrectValuesMatchingRates(t *testing.T) {
	assert := assert.New(t)
	amount := 1100

	var lenders lender.Lenders
	l1 := lender.Lender{"", 0.051, 1000}
	l2 := lender.Lender{"", 0.051, 200}

	lenders = append(lenders, l1)
	lenders = append(lenders, l2)

	q := quote{}
	q.RequestedAmount = amount
	q.lenders = lenders
	q.loanPeriodMonths = 36

	err := q.calculate()

	assert.Nil(err, "Expected error to be nil as input information was valid")

	// blended rate = (0.051*1000 + 0.051*100) / 1100
	// blended rate = 0.051
	assert.Equal("0.051", fmt.Sprintf("%.3f", q.Rate), "Expected rate to be 5.1%")

	// monthly repayment = amount * (month_rate*(1+month_rate)^num_periods)/((1+month_rate)^num_periods-1)
	// l1 = 1000 * ((0.051/12)*(1+(0.051/12))^36)/((1+(0.051/12))^36-1)
	// l1 = 30.015815509
	// l2 = 100 * ((0.051/12)*(1+(0.051/12))^36)/((1+(0.051/12))^36-1)
	// l2 = 3.0015815509
	// total monthly repayment = 33.0173970599
	assert.Equal("33.02", fmt.Sprintf("%.2f", q.MonthlyRepayment), "Expected total payment to be £33.02")

	// total repayment = 33.0173970599 * 36
	// total repayment = 1188.6262942
	assert.Equal("1188.63", fmt.Sprintf("%.2f", q.TotalRepayment), "Expected total payment to be £1188.63")
}

func TestNewQuoteCallsValidateAndFails(t *testing.T) {
	amount := MinAmount - 1

	var lenders lender.Lenders
	l1 := lender.Lender{"", 0.05, 1000}
	l2 := lender.Lender{"", 0.05, 200}

	lenders = append(lenders, l1)
	lenders = append(lenders, l2)

	quote, err := NewQuote(amount, 36, lenders)

	assert.NotNil(t, err, "Expected an error as requested amount was too low")
	assert.Nil(t, quote, "Expected a nil quote due to validation error")
}

func TestNewQuoteCallsCalculateAndFailsWhenAmountNotAvailable(t *testing.T) {
	amount := 1300

	// only 1200 available in total
	var lenders lender.Lenders
	l1 := lender.Lender{"", 0.01, 1000}
	l2 := lender.Lender{"", 0.01, 200}

	lenders = append(lenders, l1)
	lenders = append(lenders, l2)

	quote, err := NewQuote(amount, 36, lenders)

	assert.NotNil(t, err, "Expected an error as pool amount is insufficient")
	assert.Nil(t, quote, "Expected a nil quote due to validation error")
}

func TestNewQuoteCallsCalculateAndWorks(t *testing.T) {
	assert := assert.New(t)
	amount := 1200

	var lenders lender.Lenders
	l1 := lender.Lender{"", 0.051, 1000}
	l2 := lender.Lender{"", 0.069, 200}

	lenders = append(lenders, l1)
	lenders = append(lenders, l2)

	quote, err := NewQuote(amount, 36, lenders)

	assert.Nil(err, "Expected error to be nil as input information was valid")
	assert.NotNil(quote, "Expected a valid quote structure returned")

	// blended rate = (0.051*1000 + 0.069*200) / 1200
	// blended rate = 0.054
	assert.Equal("0.054", fmt.Sprintf("%.3f", quote.Rate), "Expected rate to be 5.4%")

	// monthly repayment = amount * (month_rate*(1+month_rate)^num_periods)/((1+month_rate)^num_periods-1)
	// l1 = 1000 * ((0.051/12)*(1+(0.051/12))^36)/((1+(0.051/12))^36-1)
	// l1 = 30.015815509
	// l2 = 200 * ((0.069/12)*(1+(0.069/12))^36)/((1+(0.069/12))^36-1)
	// l2 = 6.1662791696
	// total monthly repayment = 36.1820946786
	assert.Equal("36.18", fmt.Sprintf("%.2f", quote.MonthlyRepayment), "Expected total payment to be £36.18")

	// total repayment = 36.1820946786 * 36
	// total repayment = 1302.5554084
	assert.Equal("1302.56", fmt.Sprintf("%.2f", quote.TotalRepayment), "Expected total payment to be £1302.56")
}

func TestQuoteStringReturnsCorrectFormat(t *testing.T) {
	q := quote{}
	q.RequestedAmount = 1200
	q.MonthlyRepayment = 36.1820946786
	q.TotalRepayment = 1302.5554084
	q.Rate = 0.054
	resp := q.String()

	assert.NotEmpty(t, resp, "Expected a string response")

	lines := strings.Split(resp, "\n")
	assert.Equal(t, 4, len(lines), "Expected 4 files in response")
	assert.Equal(t, "Requested amount: £1200", lines[0], "Amount line should match format")
	assert.Equal(t, "Rate: 5.4%", lines[1], "Rate line should match format")
	assert.Equal(t, "Monthly repayment: £36.18", lines[2], "Monthly line should match format")
	assert.Equal(t, "Total repayment: £1302.56", lines[3], "Total line should match format")

}
