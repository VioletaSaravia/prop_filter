package lib

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
)

func Parse(query *SearchQuery) (*[]Property, error) {
	var err error
	var properties []Property

	var file *os.File
	if query.InputFile != "" {
		if file, err = os.Open(query.InputFile); err != nil {
			return nil, err
		}
	} else {
		file = os.Stdin
	}
	defer file.Close()

	if err = gocsv.UnmarshalFile(file, &properties); err == nil {
		query.OutputType = typeCSV
		return &properties, err
	}

	if err = json.NewDecoder(file).Decode(&properties); err == nil {
		query.OutputType = typeJSON
		return &properties, err
	}

	return nil, fmt.Errorf("cannot parse input data as CSV or JSON")
}

func Filter(filter SearchFilter, ctx *cli.Context) error {
	query, err := NewSearchQuery(filter, ctx)
	if err != nil {
		return err
	}

	data, err := Parse(query)
	if err != nil {
		return err
	}

	var filteredData []Property
	for _, p := range *data {
		passed := query.Args.Filter(p)
		if query.Exclude {
			passed = !passed
		}

		if passed {
			filteredData = append(filteredData, p)
		}
	}

	if err = Print(filteredData, query.OutputFile, query.OutputType); err != nil {
		return err
	}

	return nil
}

func Print(data []Property, out string, fileType FileType) (err error) {
	var file *os.File
	if out == "" {
		file = os.Stdout
	} else {
		file, err = os.Create(out)
		if err != nil {
			return err
		}
	}

	switch fileType {
	case typeCSV:
		if err = gocsv.MarshalFile(data, file); err != nil {
			return err
		}
	case typeJSON:
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		file.Write(jsonData)
	default:
		return fmt.Errorf("unreachable")
	}

	return nil
}
