package lib

import "github.com/urfave/cli/v2"

var Flags []cli.Flag = []cli.Flag{
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

func Run(ctx *cli.Context) (err error) {
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
}
