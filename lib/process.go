package lib

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/gocarina/gocsv"
)

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
		// This is much simpler at the cost of performance.

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

	switch fileType {
	case TypeCSV:

		if err = gocsv.MarshalFile(data, file); err != nil {
			return err
		}
	case TypeJSON:
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
