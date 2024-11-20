package main

type LightingType string

const (
	lightingLow  LightingType = "low"
	lightingMed  LightingType = "medium"
	lightingHigh LightingType = "high"
)

type InputType string

const (
	TypeCSV  InputType = "csv"
	TypeTSV  InputType = "tsv"
	TypeJSON InputType = "json"
)

// 2-decimal fixed
type Currency int

type Vector2 [2]float32
type Vector2D struct {
	X float32
	Y float32
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
