package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	response := make(map[string]interface{})

	if len(os.Args) != 2 {
		response["success"] = false
		response["summary"] = "Exactly one argument expected"
		response["error"] = map[string]string{
			"msg": "cvet expects exactly one argument which is the path to the csv file being vetted",
			"fix": fmt.Sprintf("call cvet with the path to the csv file as the first argument. Eg %s /path/to/csv/file", os.Args[0]),
		}
		json.NewEncoder(os.Stdout).Encode(response)
		return
	}

	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		response["success"] = false
		response["summary"] = "An internal error occured"
		response["error"] = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file.",
		}
		json.NewEncoder(os.Stdout).Encode(response)
		return
	}
	defer csvFile.Close()

	validRecords, invalidRecords, err := parse(csvFile)

	if err != nil {
		response["success"] = false
		response["summary"] = "An internal error occured"
		response["error"] = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		json.NewEncoder(os.Stdout).Encode(response)
		return
	}

	response["success"] = true
	response["summary"] = "File vetted successfully"
	response["data"] = map[string]interface{}{
		"validRecords":   validRecords,
		"invalidRecords": invalidRecords,
	}
	json.NewEncoder(os.Stdout).Encode(response)
}

// InvalidRecord is a record that is not valid.
type InvalidRecord struct {
	RowNumber int      `json:"row"`
	Columns   []string `json:"columns"`
}

// parse validates f as a valid csv file with valid data.
func parse(f io.Reader) (validRecords [][]string, invalidRecords []InvalidRecord, err error) {
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	uploadedCSV, err := r.ReadAll()
	if err != nil {
		return validRecords, invalidRecords, err
	}

	rowOffset := 2
	header := uploadedCSV[0]
	headerLength := len(header)

	for row, record := range uploadedCSV[1:] {
		currentRecord := new(InvalidRecord)
		currentRecord.RowNumber = row + rowOffset
		recordIsValid := true

		for column, field := range record {
			if strings.TrimSpace(field) == "" {
				recordIsValid = false
				currentRecord.Columns = append(currentRecord.Columns, header[column])
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
