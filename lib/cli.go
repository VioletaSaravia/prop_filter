package lib

import (
	"github.com/urfave/cli/v2"
)

var Flags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:    "input",
		Value:   "",
		Aliases: []string{"i"},
		Usage:   "Input file. Supported types: CSV, JSON.",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Output results to `FILE`.",
	},
	&cli.BoolFlag{
		Name:    "exclude",
		Aliases: []string{"e"},
		Usage:   "Exclude search",
	},
}

var Commands []*cli.Command = []*cli.Command{
	{
		Name:      "footage",
		Aliases:   []string{"f"},
		Usage:     "filter by square footage.",
		UsageText: "Property Filter footage [minimum] [maximum]",
		Action: func(ctx *cli.Context) error {
			return Filter(&FootageFilter{}, ctx)
		},
	},
	{
		Name:      "price",
		Aliases:   []string{"p"},
		Usage:     "filter by price.",
		UsageText: "Property Filter price [minimum] [maximum]",
		Args:      true,
		Action: func(ctx *cli.Context) error {
			return Filter(&PriceFilter{}, ctx)
		},
	},
	{
		Name:      "lighting",
		Aliases:   []string{"light"},
		Usage:     "filter by light level. Supported levels: min, med and max",
		UsageText: "Property Filter lighting [level]",
		Action: func(ctx *cli.Context) error {
			filter := LightingType("")
			return Filter(&filter, ctx)
		},
	},
	{
		Name:      "rooms",
		Aliases:   []string{"r"},
		Usage:     "filter by number of rooms.",
		UsageText: "Property Filter rooms [minimum] [maximum]",
		Action: func(ctx *cli.Context) error {
			return Filter(&RoomsFilter{}, ctx)
		},
	},
	{
		Name:      "bathrooms",
		Aliases:   []string{"b"},
		Usage:     "filter by number of bathrooms.",
		UsageText: "Property Filter bathrooms [minimum] [maximum]",
		Action: func(ctx *cli.Context) error {
			return Filter(&BathroomsFilter{}, ctx)
		},
	},
	{
		Name:      "location",
		Aliases:   []string{"l"},
		Usage:     "filter by distance to a location.",
		UsageText: "Property Filter location [x location] [y location] [x radius] [y radius]",
		Action: func(ctx *cli.Context) error {
			return Filter(&Location{}, ctx)
		},
	},
	{
		Name:      "description",
		Aliases:   []string{"d"},
		Usage:     "filter by description. Supports Regex.",
		UsageText: "Property Filter description [search query]",
		Action: func(ctx *cli.Context) error {
			return Filter(&DescriptionQuery{}, ctx)
		},
	},
	{
		Name:      "ammenities",
		Aliases:   []string{"a"},
		Usage:     "filter by included ammenities.",
		UsageText: "Property Filter ammenities [ammenities]",
		Action: func(ctx *cli.Context) error {
			filter := AmmenitiesFilter("")
			return Filter(&filter, ctx)
		},
	}}

var App cli.App = cli.App{
	Flags:    Flags,
	Commands: Commands,
	HelpName: "Property Filter",
	Name:     "propfilter",
	Usage:    "Filter large sets of real estate properties based on their particular attributes.",
}
