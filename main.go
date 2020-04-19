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

func sendOutput(payload jsonPayload, dest io.Writer) {
	encoder := json.NewEncoder(dest)
	if *pretty {
		encoder.SetIndent("", "   ")
	}
	encoder.Encode(payload)
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	output := run(filename)
	sendOutput(output, os.Stdout)
}

func run(filename string) (out jsonPayload) {
	out = jsonPayload{}
	if filename == "" {
		out.Ok = false
		out.Msg = "Exactly one argument expected"
		out.Error = map[string]string{
			"msg": "cvet expects exactly one argument which is the path to the csv file being vetted",
			"fix": fmt.Sprintf("call cvet with the path to the csv file as the first argument. Eg %s /path/to/csv/file", os.Args[0]),
		}
		return
	}

	csvFile, err := os.Open(filename)
	if err != nil {
		out.Ok = false
		out.Msg = fmt.Sprintf("Could not open file: %s", filename)
		out.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to open the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file",
		}
		return
	}
	defer csvFile.Close()

	// Perform the actual parsing of the csv file
	validRecords, invalidRecords, err := cval.Validate(csvFile)

	if err != nil {
		out.Ok = false
		out.Msg = "An internal error occured"
		out.Error = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		return
	}

	out.Ok = true
	out.Msg = "File vetted successfully"
	out.Data = map[string]interface{}{
		"validRecords":   validRecords,
		"invalidRecords": invalidRecords,
	}
	return
}
