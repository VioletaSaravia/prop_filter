package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
)

type FileType string

const (
	typeCSV  FileType = "csv"
	typeJSON FileType = "json"
)

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Vector2) UnmarshalCSV(s string) error {
	parts := strings.Split(s[1:len(s)-1], ",")

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("invalid X value: %w", err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("invalid Y value: %w", err)
	}

	v.X = x
	v.Y = y
	return nil
}

func (v *Vector2) MarshalCSV() (string, error) {
	return fmt.Sprintf("[%f, %f]", v.X, v.Y), nil
}

type Ammenities map[string]bool

func (a *Ammenities) UnmarshalCSV(s string) error {
	var parsed map[string]bool

	if err := json.Unmarshal([]byte(s), &parsed); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	*a = parsed
	return nil
}

func (a *Ammenities) MarshalCSV() (string, error) {
	if jsonData, err := json.Marshal(a); err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	} else {
		return string(jsonData), nil
	}
}

type Property struct {
	SquareFootage int          `json:"squareFootage" csv:"squareFootage"`
	Lighting      LightingType `json:"lighting" csv:"lighting"`
	Price         int          `json:"price" csv:"price"` // 2-decimal fixed point (cents)
	Rooms         int          `json:"rooms" csv:"rooms"`
	Bathrooms     int          `json:"bathrooms" csv:"bathrooms"`
	Location      Vector2      `json:"location" csv:"location"`
	Description   string       `json:"description" csv:"description"`
	Ammenities    Ammenities   `json:"ammenities" csv:"ammenities"`
}

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
