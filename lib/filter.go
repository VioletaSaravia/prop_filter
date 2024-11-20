package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
)

func Filter(filter SearchFilter, ctx *cli.Context) error {
	query, err := NewSearchQuery(filter, ctx)
	if err != nil {
		return err
	}

	data, err := query.Unmarshal()
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

	var ext FileType
	switch filepath.Ext(out) {
	case ".json":
		ext = typeJSON
	case ".csv":
		ext = typeCSV
	default:
		ext = fileType
	}

	switch ext {
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
