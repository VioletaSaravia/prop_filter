package lib

import (
	"github.com/urfave/cli/v2"
)

type FileType string

const (
	typeCSV  FileType = "csv"
	typeJSON FileType = "json"
)

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
