package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gocarina/gocsv"
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
		return nil, fmt.Errorf("unsupported file type: %s", ctx.String("input"))
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

func Parse(inputFile string, inputType FileType) (*[]Property, error) {
	var err error
	var properties []Property

	var file *os.File
	if inputFile != "" {
		if file, err = os.Open(inputFile); err != nil {
			return nil, err
		}

	} else {
		file = os.Stdin
	}
	defer file.Close()

	switch inputType {
	case TypeCSV:

		if err = gocsv.UnmarshalFile(file, &properties); err != nil {
			return nil, err
		}
	case TypeJSON:
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&properties); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unreachable")
	}

	return &properties, err
}

func Filter(query SearchQuery, input []Property) (output []Property, err error) {
PropertiesLoop:
	for _, p := range input {
		// NOTE: We are currently running every filter on every query, regardless of whether
		// it was passed by the user or not, by leveraging meaningful zero values in SearchQuery.
		// This is much simpler, at the cost of performance.

		if query.Description != nil && !query.Description.MatchString(p.Description) {
			continue
		}

		if p.SquareFootage < query.Footage[0] || p.SquareFootage >= query.Footage[1] {
			continue
		}

		if p.Rooms < query.Rooms[0] || p.Rooms >= query.Rooms[1] {
			continue
		}

		if p.Bathrooms < query.Bathrooms[0] || p.Bathrooms >= query.Bathrooms[1] {
			continue
		}

		if p.Price < query.Price[0] || p.Price >= query.Price[1] {
			continue
		}

		if query.Lighting != "" && p.Lighting != query.Lighting {
			continue
		}

		q := query.Location.Center
		r := query.Location.Radius
		l := p.Location
		if math.Sqrt(math.Pow(l[0]-q[0], 2)+math.Pow(l[1]-q[1], 2)) > r {
			continue
		}

		for _, a := range query.Ammenities {
			if has, ok := p.Ammenities[a]; !ok || !has {
				continue PropertiesLoop
			}
		}

		output = append(output, p)

	}
	return output, nil
}

func Print(data []Property, out string, fileType FileType) (err error) {
	if out == "" {
		fmt.Println(data)
		return nil
	}

	file, err := os.Create(out)
	if err != nil {
	return err
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Value:   "csv",
				Aliases: []string{"i"},
		Usage:   "Input type. Only necessary when reading from STDIN. Supported: CSV, JSON.",
			},
			&cli.StringFlag{
				Name:    "footage",
				Aliases: []string{"f"},
				Usage:   "Filter square footage by `RANGE`.",
			},
			&cli.StringFlag{
				Name:    "price",
				Aliases: []string{"p"},
				Usage:   "Filter price by `RANGE`.",
			},
			&cli.StringFlag{
				Name:    "lighting",
				Aliases: []string{"light"},
				Usage:   "Filter by a certain `LEVEL` of lighting. Options: low, medium, high.",
			},
			&cli.StringFlag{
				Name:    "rooms",
				Aliases: []string{"r"},
				Usage:   "Filter rooms by `RANGE`.",
			},
			&cli.StringFlag{
				Name:    "bathrooms",
				Aliases: []string{"b"},
				Usage:   "Filter bathrooms by `RANGE`.",
			},
			&cli.StringFlag{
				Name:    "location",
				Value:   "csv",
				Aliases: []string{"l"},
				Usage:   "Filter by location.",
			},
			&cli.StringFlag{
				Name:    "description",
				Aliases: []string{"d"},
				Usage:   "Filter description by `QUERY`. Supports regex queries.",
			},
			&cli.StringFlag{
				Name:    "ammenities",
				Aliases: []string{"a", "ai"},
				Usage:   "Filter by included ammenities.",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output results to `FILE`.",
			},
}

func main() {
	app := &cli.App{
		Flags: Flags,
		Name:  "Property Filter",
		Usage: "Filter large sets of real estate properties based on their particular attributes.",
		Action: func(ctx *cli.Context) (err error) {
			query, err := NewSearchQuery(ctx)
			if err != nil {
				return err
			}

			data, err := Parse(query.InputFile, query.InputType)
			if err != nil {
				return err
			}

			filteredData, err := Filter(*query, *data)
			if err != nil {
				return err
			}

			if err = Print(filteredData, query.OutputFile, query.OutputType); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
