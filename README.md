# Property Filter CLI Tool
## Overview
Property Filter is a command-line tool designed to filter large sets of real estate properties based on specific attributes, including square footage, price, number of rooms, lighting levels, and more. This tool supports CSV and JSON formats for input and output, enabling users to filter property data effectively for analysis or reporting purposes.

## Features
* Filter by any of the fields of a property (square footage, light level, price)
* Read CSV or JSON from either a file or `stdin`.
* Output results to a file or `stdout`.
* Idempotent queries that leverage `stdin`, `stdout` and unix pipes, allowing you to create complex searches by concatenating multiple queries (see examples below)

## Installation
Ensure that you have Go installed (version 1.18+). To install the tool, clone this repository and build the project:

```bash
git clone https://github.com/VioletaSaravia/property_filter
cd property_filter
go build .
```

After building, the executable `prop_filter` will be available in the current directory. 

## Documentation

You can access detailed documentation by running `./prop_filter help`. See also `prop_filter [command] help` for additional information on the parameters of each command.

```text
NAME:
   Property Filter - Filter large sets of real estate properties based on their particular attributes.

USAGE:
   Property Filter [global options] command [command options]

COMMANDS:
   footage, f       filter by square footage.
   price, p         filter by price.
   lighting, light  filter by light level. Supported levels: min, med and max
   rooms, r         filter by number of rooms.
   bathrooms, b     filter by number of bathrooms.
   location, l      filter by distance to a location.
   description, d   filter by description. Supports Regex.
   ammenities, a    filter by included ammenities.
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value  Input file. Supported types: CSV, JSON.
   --output FILE, -o FILE   Output results to FILE.
   --exclude, -e            Exclude search (default: false)
   --help, -h               show help
```

## Examples

Simple queries can be performed by passing the desired field to search, plus its arguments and any optional parameters (`--input`, `--output` and `--exclude` to exclude the search results)

```bash
# Searches for properties with 4+ rooms in examples/data.csv, save to big.csv
./prop_filter -i examples/data.csv -o big.csv rooms 4

# Searches for properties without low light in examples/data.json
./prop_filter -e -i examples/data.csv light low

# Searches for properties within 100 units of (0, 0) in examples/data.csv, save to search.json
./prop_filter -i examples/data.csv -o search.json location 0 0 100

# Searches for properties whose description matches the regex "[Ss]pa..ous" in examples/data.csv
./prop_filter -i examples/data.csv description "[Ss]pa..ous"
```

By piping query results to `stdout` and taking them from `stdin`, it's possible to assemble complex searches in a single line:

```bash
# Search for properties within 100 units of (0, 0), with more than four rooms and a pool, save to search.csv
./prop_filter -i examples/data.csv location 0 0 100 | ./prop_filter rooms 4 | ./prop_filter -o search.csv a pool

# Search for properties under 1600 sqft, with a gym and parking
./prop_filter -i examples/data.json -e footage 1600 | ./prop_filter a gym | ./prop_filter a parking
```