package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yeboahnanaosei/gitplus/csv"
)

func main() {
	out := make(map[string]interface{})

	csvFilepath := flag.String("file", "", "The path to the file the CSV file")
	flag.Parse()

	if *csvFilepath == "" {
		out["success"] = false
		out["code"] = 400
		out["summary"] = "Invalid request"
		out["error"] = map[string]string{
			"msg": "No path to csv file supplied. A path to a csv file was expected but none was supplied",
			"fix": "use the '-file' flag to supply the path eg. -file /path/to/csv/file",
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	csvFile, err := os.Open(*csvFilepath)
	if err != nil {
		out["success"] = false
		out["code"] = 500
		out["summary"] = "An internal error occured"
		out["error"] = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file. Server said: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}
	defer csvFile.Close()

	validRecords, invalidRecords, err := csv.Validate(csvFile)

	if err != nil {
		out["success"] = false
		out["code"] = 500
		out["summary"] = "An internal error occured"
		out["error"] = map[string]string{
			"msg": fmt.Sprintf("There was an error trying to process the csv file. Server said: %v", err),
			"fix": "Ensure you provided a valid csv file. If this continues, please wait and try again later. You can also contact support",
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	out["success"] = true
	out["code"] = 200
	out["summary"] = "Operation successful"
	out["data"] = map[string]interface{}{
		"validRecords":   validRecords,
		"invalidRecords": invalidRecords,
	}
	json.NewEncoder(os.Stdout).Encode(out)
	return
}
