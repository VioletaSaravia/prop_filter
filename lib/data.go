package lib

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v2"
)

type LightingType string

const (
	lightingLow  LightingType = "low"
	lightingMed  LightingType = "medium"
	lightingHigh LightingType = "high"
)

type FileType string

const (
	TypeCSV  FileType = "csv"
	TypeJSON FileType = "json"
)

// 2-decimal fixed
type Currency int

type Vector2 [2]float64

func (v *Vector2) UnmarshalCSV(string) error {
	return nil
}

type Ammenities map[string]bool

func (a *Ammenities) UnmarshalCSV(string) error {
	return nil
}

type Property struct {
	SquareFootage int          `json:"squareFootage" csv:"squareFootage"`
	Lighting      LightingType `json:"lighting" csv:"lighting"`
	Price         Currency     `json:"price" csv:"price"`
	Rooms         int          `json:"rooms" csv:"rooms"`
	Bathrooms     int          `json:"bathrooms" csv:"bathrooms"`
	Location      Vector2      `json:"location" csv:"location"`
	Description   string       `json:"description" csv:"description"`
	Ammenities    Ammenities   `json:"ammenities" csv:"ammenities"`
}

type Area struct {
	Center Vector2
	Radius float64
}

type SearchQuery struct {
	InputFile   string
	OutputFile  string
	InputType   FileType
	OutputType  FileType
	Footage     [2]int
	Lighting    LightingType
	Rooms       [2]int
	Bathrooms   [2]int
	Price       [2]Currency
	Location    Area
	Description *regexp.Regexp
	Ammenities  []string
}

func NewSearchQuery(ctx *cli.Context) (query *SearchQuery, err error) {
	query = &SearchQuery{}

	// in file
	query.InputFile = ctx.Args().Get(0)

	// light
	switch ctx.String("lighting") {
	case "low":
		query.Lighting = lightingLow
	case "medium", "med":
		query.Lighting = lightingMed
	case "high":
		query.Lighting = lightingHigh
	case "":
		break
	default:
		return nil, fmt.Errorf("unsupported lighting level: %s", ctx.String("lighting"))
	}

	// type
	if query.InputFile == "" && ctx.String("input") == "" {
		return nil, fmt.Errorf("input type must be specified with -i when reading from stdin")
	}

	var ext string
	if query.InputFile != "" {
		ext = filepath.Ext(query.InputFile)
	} else {
		ext = ctx.String("input")
	}

	switch ext {
	case "csv", ".csv":
		query.InputType = TypeCSV
	case "json", ".json":
		query.InputType = TypeJSON
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	query.OutputFile = ctx.String("output")

	switch filepath.Ext(query.OutputFile) {
	case ".csv":
		query.OutputType = TypeCSV
	case ".json":
		query.OutputType = TypeJSON
	case "":
		query.OutputType = query.InputType
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ctx.String("output"))
	}

	// rooms
	if ctx.String("rooms") == "" {
		query.Rooms = [2]int{0, math.MaxInt}
	}

	// bathrooms
	if ctx.String("bathrooms") == "" {
		query.Bathrooms = [2]int{0, math.MaxInt}
	}

	// loc
	if ctx.String("location") == "" {
		query.Location = Area{
			Center: Vector2{0, 0},
			Radius: math.MaxFloat32,
		}
	}

	// desc
	if query.Description, err = regexp.Compile(ctx.String("description")); err != nil {
		return nil, err
	}

	// amm
	// TODO

	// price
	if ctx.String("price") == "" {
		query.Price = [2]Currency{0, math.MaxInt}
	}

	// footage
	if ctx.String("footage") == "" {
		query.Footage = [2]int{0, math.MaxInt}
	}

	return query, err
}
