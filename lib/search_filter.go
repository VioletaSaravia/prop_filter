package lib

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/urfave/cli/v2"
)

type SearchFilter interface {
	Parse(cli.Args) error
	Filter(data Property) bool
}

type LightingType string

const (
	lightingLow  LightingType = "low"
	lightingMed  LightingType = "medium"
	lightingHigh LightingType = "high"
	lightingNone LightingType = ""
)

func (t *LightingType) Parse(args cli.Args) error {
	val := LightingType(args.Get(0))

	switch val {
	case lightingHigh, lightingLow, lightingMed, lightingNone:
		*t = val
	default:
		return fmt.Errorf("invalid light type: %s", val)
	}

	return nil
}

func (t LightingType) Filter(data Property) bool {
	return t == LightingType("") || data.Lighting == t
}

type IntRange [2]int

func (i *IntRange) Parse(args cli.Args) error {
	min, err := strconv.Atoi(args.Get(0))
	if err != nil {
		return fmt.Errorf("invalid argument: %s", args.Get(0))
	}

	var max int = math.MaxInt
	if args.Len() >= 2 {
		max, err = strconv.Atoi(args.Get(1))
		if err != nil {
			return fmt.Errorf("invalid argument: %s", args.Get(1))
		}
	}
	i[0] = min
	i[1] = max

	return nil
}

type RoomsFilter struct{ IntRange }

func (f RoomsFilter) Filter(data Property) bool {
	return data.Rooms >= f.IntRange[0] && data.Rooms <= f.IntRange[1]
}

type PriceFilter struct{ IntRange }

func (f PriceFilter) Filter(data Property) bool {
	return data.Price >= f.IntRange[0] && data.Price <= f.IntRange[1]
}

type BathroomsFilter struct{ IntRange }

func (f BathroomsFilter) Filter(data Property) bool {
	return data.Bathrooms >= f.IntRange[0] && data.Bathrooms <= f.IntRange[1]
}

type FootageFilter struct{ IntRange }

func (f FootageFilter) Filter(data Property) bool {
	return data.SquareFootage >= f.IntRange[0] && data.SquareFootage <= f.IntRange[1]
}

type Location struct {
	Center Vector2
	Radius float64
}

func (l *Location) Parse(args cli.Args) error {
	if args.Len() < 3 {
		return fmt.Errorf("insufficient arguments")
	}

	x, err := strconv.ParseFloat(args.Get(0), 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(args.Get(1), 64)
	if err != nil {
		return err
	}
	r, err := strconv.ParseFloat(args.Get(2), 64)
	if err != nil {
		return err
	}

	l.Center = Vector2{x, y}
	l.Radius = r

	return nil
}

func (l Location) Filter(data Property) bool {
	q := l.Center
	r := l.Radius
	loc := data.Location
	return math.Sqrt(math.Pow(loc[0]-q[0], 2)+math.Pow(loc[1]-q[1], 2)) <= r
}

type DescriptionQuery struct{ *regexp.Regexp }

func (i *DescriptionQuery) Parse(args cli.Args) (err error) {
	i.Regexp, err = regexp.Compile(args.Get(0))
	return err
}

func (d DescriptionQuery) Filter(data Property) bool {
	return d.MatchString(data.Description)
}

type AmmenitiesFilter string

func (i *AmmenitiesFilter) Parse(args cli.Args) (err error) {
	*i = AmmenitiesFilter(args.Get(0))
	return nil
}

func (a AmmenitiesFilter) Filter(data Property) bool {
	has, ok := data.Ammenities[string(a)]
	return ok && has
}
