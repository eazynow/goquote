package lender

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVImportFailsOnNotFound(t *testing.T) {
	filename := "unknownfile.csv"

	_, err := ImportCSV(filename)

	assert.NotNil(t, err, "Expected an error as the csv file cannot be found")
}

func TestCSVImportFailsWithMissingFields(t *testing.T) {
	filename := "test_missing_data.csv"
	_, err := ImportCSV(filename)

	assert.NotNil(t, err, "Expected an error as the csv was missing a field in row 2")
}

func TestCSVImportFailsWithBadRate(t *testing.T) {
	filename := "test_bad_rate.csv"
	_, err := ImportCSV(filename)

	if err == nil {
		t.Error("Expected an error as the csv has bad data in the rate field")
	}

	assert.Equal(
		t,
		"Error unmarshalling rate field on line 4. Cause: strconv.ParseFloat: parsing \"bad_rate\": invalid syntax",
		err.Error(),
		"The error message was not formatted correctly")
}

func TestCSVImportFailsWithBadAmount(t *testing.T) {
	filename := "test_bad_amount.csv"
	_, err := ImportCSV(filename)

	assert.NotNil(t, err, "Expected an error as the csv has bad data in the rate field")
}

func TestCSVImportWorksWithGoodData(t *testing.T) {
	filename := "test_good_test.csv"

	lenders, err := ImportCSV(filename)

	assert.Nil(t, err, "Expected file to work but it failed")
	assert.NotNil(t, lenders, "Expected to get 2 lenders back but did not get any")

	assert.Equal(t, 2, len(lenders), "There should be 2 loenders")

	assert.Equal(t, lenders[0].Name, "Lender1")
	assert.Equal(t, lenders[0].Rate, 0.075)
	assert.Equal(t, lenders[0].Available, 640)

	assert.Equal(t, lenders[1].Name, "Lender2")
	assert.Equal(t, lenders[1].Rate, 0.069)
	assert.Equal(t, lenders[1].Available, 480)
}
