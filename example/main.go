package main

import (
	"encoding/json"
	"fmt"

	"github.com/bclswl0827/travel"
)

func main() {
	table, err := travel.NewAK135()
	if err != nil {
		panic(err)
	}

	res := table.Estimate(travel.GetDeltaByCoordinates(
		// earthquake
		28.25, 104.87,
		// station
		29.81, 106.40,
	), 10, true)
	data, _ := json.MarshalIndent(res, "", "  ")

	fmt.Println(string(data))
}
