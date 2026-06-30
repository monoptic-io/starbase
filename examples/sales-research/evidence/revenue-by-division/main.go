// Recomputes total 2025 revenue by division, descending.
//
//starbase:deps ../../data/sales.csv
package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("../../data/sales.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rows, _ := csv.NewReader(f).ReadAll()
	col := map[string]int{}
	for i, h := range rows[0] {
		col[h] = i
	}
	rev := map[string]int{}
	for _, r := range rows[1:] {
		v, _ := strconv.Atoi(r[col["revenue"]])
		rev[r[col["division"]]] += v
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
	table := [][]string{{"division", "total"}}
	for _, e := range list {
		table = append(table, []string{e.d, strconv.Itoa(e.t)})
	}
	json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"table": table})
}
