package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yeboahnanaosei/go/cval"
)

type jsonResponse struct {
	Success bool              `json:"success"`
	Summary string            `json:"summary"`
	Data    interface{}       `json:"data,omitempty"`
	Error   map[string]string `json:"error,omitempty"`
}

func main() {
	payload := jsonResponse{}

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
	validRecords, invalidRecords, err := cval.Validate(csvFile)

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
