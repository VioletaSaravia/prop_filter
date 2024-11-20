package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
)

type SearchQuery struct {
	InputFile  string
	OutputFile string
	OutputType FileType
	Exclude    bool
	Args       SearchFilter
}

func NewSearchQuery(filter SearchFilter, ctx *cli.Context) (query *SearchQuery, err error) {
	query = &SearchQuery{
		InputFile:  ctx.String("input"),
		OutputFile: ctx.String("output"),
		OutputType: FileType(""),
		Exclude:    ctx.Bool("exclude"),
		Args:       filter,
	}

	if err = query.Args.Parse(ctx.Args()); err != nil {
		return nil, err
	}

	return query, err
}

func (query *SearchQuery) Unmarshal() (*[]Property, error) {
	var err error
	var data []byte
	var properties []Property

	if query.InputFile != "" {
		if data, err = os.ReadFile(query.InputFile); err != nil {
			return nil, err
		}
	} else {
		if data, err = io.ReadAll(os.Stdin); err != nil {
			return nil, err
		}
	}

	var csvErr error
	if csvErr = gocsv.UnmarshalBytes(data, &properties); csvErr == nil {
		query.OutputType = typeCSV
		return &properties, nil
	}

	var jsonErr error
	if jsonErr = json.Unmarshal(data, &properties); jsonErr == nil {
		query.OutputType = typeJSON
		return &properties, nil
	}

	return nil, fmt.Errorf("cannot parse input data as CSV (%s) or JSON (%s)", csvErr, jsonErr)
}

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
