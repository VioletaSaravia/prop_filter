package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		input  string
		output string
	}{{
		name:  "Rooms filter/single argument/CSV input",
		args:  []string{os.Args[0], "-e", "-i", "examples/data.csv", "rooms", "4"},
		input: "",
		output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
1200,medium,250000,3,2,"[37.774900, -122.419400]",A cozy family home in a quiet neighborhood.,"{""garage"":true,""pool"":false,""yard"":true}"
800,low,150000,2,1,"[40.712800, -74.006000]",Compact apartment in the city center.,"{""garage"":false,""pool"":false,""yard"":false}"
900,low,175000,2,1,"[41.878100, -87.629800]",Affordable condo with easy access to amenities.,"{""garage"":false,""pool"":false,""yard"":false}"
1600,medium,300000,3,2,"[47.606200, -122.332100]",Charming house near downtown.,"{""garage"":false,""pool"":false,""yard"":true}"
1100,low,200000,3,2,"[32.776700, -96.797000]",Charming starter home.,"{""garage"":true,""pool"":false,""yard"":true}"
1700,medium,325000,3,2,"[29.760400, -95.369800]",Contemporary design in a growing community.,"{""garage"":true,""pool"":false,""yard"":true}"
750,low,140000,1,1,"[39.952600, -75.165200]",Small but functional studio.,"{""garage"":false,""pool"":false,""yard"":false}"
1300,high,275000,3,2,"[44.977800, -93.265000]",Stylish townhouse in a prime location.,"{""garage"":true,""pool"":false,""yard"":false}"
1400,medium,290000,3,2,"[38.907200, -77.036900]",Classic home in the nation's capital.,"{""garage"":true,""pool"":false,""yard"":true}"
1150,low,210000,2,2,"[37.338200, -121.886300]",Affordable duplex in a growing tech hub.,"{""garage"":false,""pool"":false,""yard"":false}"
950,low,180000,2,1,"[42.360100, -71.058900]",Historic condo with a touch of charm.,"{""garage"":false,""pool"":false,""yard"":false}"
1450,medium,310000,3,2,"[40.712800, -74.006000]",Stylish family home in the city.,"{""garage"":true,""pool"":false,""yard"":true}"
1250,medium,265000,3,2,"[37.774900, -122.419400]",Eco-friendly design with solar panels.,"{""garage"":false,""pool"":false,""yard"":true}"
`,
	},
		{
			name:  "Bathrooms filter/two arguments/CSV input",
			args:  []string{os.Args[0], "-e", "-i", "examples/data.csv", "bathrooms", "2", "3"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
800,low,150000,2,1,"[40.712800, -74.006000]",Compact apartment in the city center.,"{""garage"":false,""pool"":false,""yard"":false}"
900,low,175000,2,1,"[41.878100, -87.629800]",Affordable condo with easy access to amenities.,"{""garage"":false,""pool"":false,""yard"":false}"
3000,high,500000,6,4,"[25.761700, -80.191800]",Luxury mansion with private pool.,"{""garage"":true,""pool"":true,""yard"":true}"
750,low,140000,1,1,"[39.952600, -75.165200]",Small but functional studio.,"{""garage"":false,""pool"":false,""yard"":false}"
3200,high,600000,5,4,"[28.538300, -81.379200]",Expansive property with luxury finishes.,"{""garage"":true,""pool"":true,""yard"":true}"
950,low,180000,2,1,"[42.360100, -71.058900]",Historic condo with a touch of charm.,"{""garage"":false,""pool"":false,""yard"":false}"
`,
		},
		{
			name:  "Location filter/CSV input",
			args:  []string{os.Args[0], "-i", "examples/data.csv", "location", "0.0", "0.0", "100.0"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
800,low,150000,2,1,"[40.712800, -74.006000]",Compact apartment in the city center.,"{""garage"":false,""pool"":false,""yard"":false}"
1800,medium,350000,4,2,"[36.162700, -86.781600]",Modern home in a bustling suburb.,"{""garage"":true,""pool"":false,""yard"":true}"
900,low,175000,2,1,"[41.878100, -87.629800]",Affordable condo with easy access to amenities.,"{""garage"":false,""pool"":false,""yard"":false}"
3000,high,500000,6,4,"[25.761700, -80.191800]",Luxury mansion with private pool.,"{""garage"":true,""pool"":true,""yard"":true}"
1700,medium,325000,3,2,"[29.760400, -95.369800]",Contemporary design in a growing community.,"{""garage"":true,""pool"":false,""yard"":true}"
750,low,140000,1,1,"[39.952600, -75.165200]",Small but functional studio.,"{""garage"":false,""pool"":false,""yard"":false}"
1400,medium,290000,3,2,"[38.907200, -77.036900]",Classic home in the nation's capital.,"{""garage"":true,""pool"":false,""yard"":true}"
3200,high,600000,5,4,"[28.538300, -81.379200]",Expansive property with luxury finishes.,"{""garage"":true,""pool"":true,""yard"":true}"
950,low,180000,2,1,"[42.360100, -71.058900]",Historic condo with a touch of charm.,"{""garage"":false,""pool"":false,""yard"":false}"
1450,medium,310000,3,2,"[40.712800, -74.006000]",Stylish family home in the city.,"{""garage"":true,""pool"":false,""yard"":true}"
`,
		},
		{
			name:  "Description filter/CSV input",
			args:  []string{os.Args[0], "-i", "examples/data.csv", "description", "desi.n"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
1700,medium,325000,3,2,"[29.760400, -95.369800]",Contemporary design in a growing community.,"{""garage"":true,""pool"":false,""yard"":true}"
1250,medium,265000,3,2,"[37.774900, -122.419400]",Eco-friendly design with solar panels.,"{""garage"":false,""pool"":false,""yard"":true}"
`,
		},
		{
			name:  "Ammenity filter/CSV input",
			args:  []string{os.Args[0], "-e", "-i", "examples/data.csv", "ammenities", "yard"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
800,low,150000,2,1,"[40.712800, -74.006000]",Compact apartment in the city center.,"{""garage"":false,""pool"":false,""yard"":false}"
900,low,175000,2,1,"[41.878100, -87.629800]",Affordable condo with easy access to amenities.,"{""garage"":false,""pool"":false,""yard"":false}"
750,low,140000,1,1,"[39.952600, -75.165200]",Small but functional studio.,"{""garage"":false,""pool"":false,""yard"":false}"
1300,high,275000,3,2,"[44.977800, -93.265000]",Stylish townhouse in a prime location.,"{""garage"":true,""pool"":false,""yard"":false}"
1150,low,210000,2,2,"[37.338200, -121.886300]",Affordable duplex in a growing tech hub.,"{""garage"":false,""pool"":false,""yard"":false}"
950,low,180000,2,1,"[42.360100, -71.058900]",Historic condo with a touch of charm.,"{""garage"":false,""pool"":false,""yard"":false}"
`,
		},
		{
			name:  "Light filter/CSV input",
			args:  []string{os.Args[0], "-i", "examples/data.csv", "lighting", "low"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
800,low,150000,2,1,"[40.712800, -74.006000]",Compact apartment in the city center.,"{""garage"":false,""pool"":false,""yard"":false}"
900,low,175000,2,1,"[41.878100, -87.629800]",Affordable condo with easy access to amenities.,"{""garage"":false,""pool"":false,""yard"":false}"
1100,low,200000,3,2,"[32.776700, -96.797000]",Charming starter home.,"{""garage"":true,""pool"":false,""yard"":true}"
750,low,140000,1,1,"[39.952600, -75.165200]",Small but functional studio.,"{""garage"":false,""pool"":false,""yard"":false}"
1150,low,210000,2,2,"[37.338200, -121.886300]",Affordable duplex in a growing tech hub.,"{""garage"":false,""pool"":false,""yard"":false}"
950,low,180000,2,1,"[42.360100, -71.058900]",Historic condo with a touch of charm.,"{""garage"":false,""pool"":false,""yard"":false}"
`,
		},
		{
			name:  "Footage filter/CSV input",
			args:  []string{os.Args[0], "-i", "examples/data.csv", "footage", "1500", "2500"},
			input: "",
			output: `squareFootage,lighting,price,rooms,bathrooms,location,description,ammenities
2500,high,450000,5,3,"[34.052200, -118.243700]",Spacious villa with scenic views.,"{""garage"":true,""pool"":true,""yard"":true}"
1800,medium,350000,4,2,"[36.162700, -86.781600]",Modern home in a bustling suburb.,"{""garage"":true,""pool"":false,""yard"":true}"
1600,medium,300000,3,2,"[47.606200, -122.332100]",Charming house near downtown.,"{""garage"":false,""pool"":false,""yard"":true}"
2200,medium,400000,4,3,"[39.739200, -104.990300]",Beautiful home with mountain views.,"{""garage"":true,""pool"":false,""yard"":true}"
1700,medium,325000,3,2,"[29.760400, -95.369800]",Contemporary design in a growing community.,"{""garage"":true,""pool"":false,""yard"":true}"
2000,medium,375000,4,3,"[33.448400, -112.074000]",Spacious house with modern upgrades.,"{""garage"":true,""pool"":true,""yard"":true}"
2100,high,425000,4,3,"[34.052200, -118.243700]",Modern and open-concept home.,"{""garage"":true,""pool"":true,""yard"":true}"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			r, w, _ := os.Pipe()
			os.Stdout = w

			main()

			_ = w.Close()

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			_ = r.Close()

			fmt.Println(buf.String())
			fmt.Println(tt.output)
			assert.Equal(t, tt.output, buf.String())
		})
	}

}
