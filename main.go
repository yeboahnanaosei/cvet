package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Summary string `json:"summary"`
	Data interface{} `json:"data,omitempty"`
	Error map[string]string `json:"error,omitempty"`
}

func main() {
	payload := new(jsonResponse)

	// Make sure this program is called with exactly one argument
	if len(os.Args) != 2 {
		payload.Success = false
		payload.Summary = "Exactly one argument expected"
		payload.Error = map[string]string{
			"msg": "cvet expects exactly one argument which is the path to the csv file being vetted",
			"fix": fmt.Sprintf("call cvet with the path to the csv file as the first argument. Eg %s /path/to/csv/file", os.Args[0]),
		}
		json.NewEncoder(os.Stdout).Encode(payload)
		return
	}

	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		payload.Success = false
		payload.Summary = "An internal error occured"
		payload.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to open the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file.",
		}
		json.NewEncoder(os.Stdout).Encode(payload)
		return
	}
	defer csvFile.Close()

	// Perform the actual parsing of the csv file
	validRecords, invalidRecords, err := parse(csvFile)

	if err != nil {
		payload.Success = false
		payload.Summary = "An internal error occured"
		payload.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		json.NewEncoder(os.Stdout).Encode(payload)
		return
	}

	payload.Success = true
	payload.Summary = "File vetted successfully"
	payload.Data = map[string]interface{}{
		"validRecords":   validRecords,
		"invalidRecords": invalidRecords,
	}
	json.NewEncoder(os.Stdout).Encode(payload)
}

// invalidRecord is a record or row in the csv file that has at least
// one empty column.
type invalidRecord struct {
	RowNumber int      `json:"row"`
	Columns   []string `json:"cols"`
}

// parse validates f as a valid csv file with valid data.
func parse(f io.Reader) (validRecords [][]string, invalidRecords []invalidRecord, err error) {
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	uploadedCSV, err := r.ReadAll()
	if err != nil {
		return validRecords, invalidRecords, err
	}

	// The first row in the csv is usually the header. Which has the name of each
	// column in the csv file
	var header []string = uploadedCSV[0]
	headerLength := len(header)

	// To determine the row number of an invalid row, we need to account for
	// the header in the file
	const headerOffset = 2

	// Skip the header. Go through each row in the file checking that for each
	// row there are no empty columns
	for rowIndex, record := range uploadedCSV[1:] {
		currentRecord := new(invalidRecord)
		currentRecord.RowNumber = rowIndex + headerOffset
		recordIsValid := true

		for columnIndex, field := range record {
			if strings.TrimSpace(field) == "" {
				recordIsValid = false
				currentRecord.Columns = append(currentRecord.Columns, header[columnIndex])
			}
		}

		if recordIsValid {
			validRecords = append(validRecords, record)
		} else if !recordIsValid && len(currentRecord.Columns) != headerLength {
			invalidRecords = append(invalidRecords, *currentRecord)
		}
	}
	return validRecords, invalidRecords, nil
}
