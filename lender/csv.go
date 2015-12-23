package lender

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

const (
	// csvFieldsPerRow is the number of expected fields on every row of the
	// csv file.
	csvFieldsPerRow = 3
	// csvNameIndex is the index reference for the Name field
	csvNameIndex = 0
	// csvRateIndex is the index reference for the Rate field
	csvRateIndex = 1
	// csvAmountIndex is the index reference for the Amount field
	csvAmountIndex = 2
)

// FieldParseError is an error structure that is used when importing a csv
// file. A custom structure is used to provide additional information
// when debugging an import problem of a csv file.
type FieldParseError struct {
	// Cause is the root cause of the parsing error.
	Cause error
	// LineNo is the line (row) that the error occurred on.
	LineNo int
	// Field is the name of the field where the error occurrd.
	Field string
}

// The Error function is used to satisfy the error interface.
// Returns a string represntation of the FieldParseError.
func (l *FieldParseError) Error() string {
	return fmt.Sprintf("Error unmarshalling %s field on line %d. Cause: %s",
		l.Field,
		l.LineNo,
		l.Cause.Error())
}

// NewFieldParseError is a helper function for creating a field parse error.
// Returns a pointer to a new FieldParseError structure.
func NewFieldParseError(lineNo int, field string, cause error) *FieldParseError {
	return &FieldParseError{
		LineNo: lineNo,
		Field:  field,
		Cause:  cause,
	}
}

// ImportCSV is a function used to import a csv file located by filename, parse
// and convert to a Lenders slice of Lender structures.
// Returns a Lenders slice of lenders if successfully imported, or
// the assoicated error otherwise.
func ImportCSV(filename string) (Lenders, error) {
	var lenders Lenders

	// Attempt to open the csv file
	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer csvfile.Close()

	// Attempt to parse the file as a csv
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = csvFieldsPerRow
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	lineNo := 1

	// Loop through all rows in the csv file.
	for _, record := range rawCSVdata {

		// Ignore the first line which contains column headers.
		if lineNo > 1 {
			l := Lender{}
			l.Name = record[csvNameIndex]

			// The rate should be a floating point number.
			l.Rate, err = strconv.ParseFloat(record[csvRateIndex], 64)
			if err != nil {
				return nil, NewFieldParseError(lineNo, "rate", err)
			}

			// The amount availale should be an integer
			l.Available, err = strconv.Atoi(record[csvAmountIndex])
			if err != nil {
				return nil, NewFieldParseError(lineNo, "available", err)
			}

			lenders = append(lenders, l)
		}
		lineNo++
	}

	return lenders, nil
}
