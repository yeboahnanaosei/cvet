package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yeboahnanaosei/go/cval"
)

type jsonPayload struct {
	Ok    bool              `json:"ok"`
	Msg   string            `json:"msg"`
	Data  interface{}       `json:"data,omitempty"`
	Error map[string]string `json:"error,omitempty"`
}

var pretty = flag.Bool("p", false, "Pretty print output")

func sendOutput(payload *jsonPayload, dest io.Writer) {
	encoder := json.NewEncoder(dest)
	if *pretty {
		encoder.SetIndent("", "   ")
	}
	encoder.Encode(payload)
}

func main() {
	flag.Parse()
	payload := jsonPayload{}
	defer sendOutput(&payload, os.Stdout)

	// Make sure this program is called with exactly one argument
	if len(os.Args) != 2 {
		payload.Ok = false
		payload.Msg = "Exactly one argument expected"
		payload.Error = map[string]string{
			"msg": "cvet expects exactly one argument which is the path to the csv file being vetted",
			"fix": fmt.Sprintf("call cvet with the path to the csv file as the first argument. Eg %s /path/to/csv/file", os.Args[0]),
		}
		return
	}

	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		payload.Ok = false
		payload.Msg = "An internal error occured"
		payload.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to open the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file.",
		}
		return
	}
	defer csvFile.Close()

	// Perform the actual parsing of the csv file
	validRecords, invalidRecords, err := cval.Validate(csvFile)

	if err != nil {
		payload.Ok = false
		payload.Msg = "An internal error occured"
		payload.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		return
	}

	payload.Ok = true
	payload.Msg = "File vetted successfully"
	payload.Data = map[string]interface{}{
		"validRecords":   validRecords,
		"invalidRecords": invalidRecords,
	}
}
