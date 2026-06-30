// Recomputes the number of Midwest sales regions. One evidence unit: re-run only
// when this code or its declared data dependency changes.
//
//starbase:deps ../../data/sales.csv
package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("../../data/sales.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rows, _ := csv.NewReader(f).ReadAll()
	div := 0
	for i, h := range rows[0] {
		if h == "division" {
			div = i
		}
	}
	n := 0
	for _, r := range rows[1:] {
		if r[div] == "Midwest" {
			n++
		}
	}
	json.NewEncoder(os.Stdout).Encode(map[string]string{"value": strconv.Itoa(n)})
}

