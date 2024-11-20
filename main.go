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
	SquareFootage int
	Lighting      LightingType
	Price         Currency
	Rooms         int
	Bathrooms     int
	Location      Vector2
	Description   string
	Ammenities    map[string]bool //yard, garage, pool, etc
}

type Area struct {
	Center Vector2D
	Radius Vector2D
}

type SearchQuery struct {
	InputFile   io.Reader
	OutputFile  io.Writer
	Type        InputType
	Footage     [2]int
	Lighting    LightingType
	Rooms       [2]int
	Bathrooms   [2]int
	Location    Area
	Description *regexp.Regexp
	Ammenities  string
}

func NewSearchQuery(ctx *cli.Context) (query *SearchQuery, err error) {
	query = &SearchQuery{}

	// in file
	if ctx.NArg() > 0 {
		if query.InputFile, err = os.Open(ctx.Args().Get(0)); err != nil {
			return nil, err
		}
	}

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
	switch ctx.String("input") {
	case "csv", "":
		query.Type = TypeCSV
	case "tsv":
		query.Type = TypeTSV
	case "json":
		query.Type = TypeJSON
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ctx.String("input"))
	}

	// rooms

	// bathrooms

	// loc

	// desc
	if query.Description, err = regexp.Compile(ctx.String("description")); err != nil {
		return nil, err
	}

	// amm

	// price

	// footage

	// output
	if outputFile := ctx.String("output"); outputFile != "" {
		if query.OutputFile, err = os.Create(outputFile); err != nil {
			return nil, err
		}
	}

	return query, err
}

func FilterAndPrint(query SearchQuery) (err error) {
	reader := bufio.NewReader(query.InputFile)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	return err
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Value:   "csv",
				Aliases: []string{"i"},
				Usage:   "Input type. Supported: CSV, TSV, JSON.",
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
		},
		Name:  "Property Filter",
		Usage: "Filter large sets of real estate properties based on their particular attributes.",
		Action: func(ctx *cli.Context) error {
			query, err := NewSearchQuery(ctx)
			if err != nil {
				return err
			}

			if err := FilterAndPrint(*query); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
