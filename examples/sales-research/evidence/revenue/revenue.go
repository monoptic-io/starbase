// Package revenue recomputes 2025 revenue by division.
package revenue

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

//starbase:deps data/sales.csv

// RevenueByDivision returns total 2025 revenue per division, descending. It
// becomes the check "revenue-by-division".
func RevenueByDivision() ([][]string, error) {
	f, err := os.Open("data/sales.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	c := map[string]int{}
	for i, h := range rows[0] {
		c[h] = i
	}
	rev := map[string]int{}
	for _, r := range rows[1:] {
		v, _ := strconv.Atoi(r[c["revenue"]])
		rev[r[c["division"]]] += v
	}
	type kv struct {
		d string
		t int
	}
	var list []kv
	for d, t := range rev {
		list = append(list, kv{d, t})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].t > list[j].t })
	out := [][]string{{"division", "total"}}
	for _, e := range list {
		out = append(out, []string{e.d, strconv.Itoa(e.t)})
	}
	return out, nil
}

