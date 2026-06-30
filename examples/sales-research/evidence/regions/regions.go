// Package regions recomputes Midwest sales-region facts. Plain functions —
// starbase discovers them and generates the runner, like `go test`.
package regions

import (
	"encoding/csv"
	"os"
)

//starbase:deps data/sales.csv

// MidwestRegions counts the distinct sales regions in the Midwest division.
// It becomes the check "midwest-regions".
func MidwestRegions() (int, error) {
	rows, err := readCSV()
	if err != nil {
		return 0, err
	}
	div := col(rows[0], "division")
	n := 0
	for _, r := range rows[1:] {
		if r[div] == "Midwest" {
			n++
		}
	}
	return n, nil
}

func readCSV() ([][]string, error) {
	f, err := os.Open("data/sales.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return csv.NewReader(f).ReadAll()
}

func col(header []string, name string) int {
	for i, h := range header {
		if h == name {
			return i
		}
	}
	return -1
}

