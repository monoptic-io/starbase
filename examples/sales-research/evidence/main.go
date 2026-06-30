// Evidence program: computes the facts the findings assert, printing them as
// JSON so `starbase verify` can diff them against the article. Pure stdlib here,
// but it could just as well shell out to duckdb, hit a DB driver, etc.
package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sort"
	"strconv"
)

type check struct {
	Value string     `json:"value,omitempty"`
	Table [][]string `json:"table,omitempty"`
}

func main() {
	f, err := os.Open("../data/sales.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rows, _ := csv.NewReader(f).ReadAll()
	col := map[string]int{}
	for i, h := range rows[0] {
		col[h] = i
	}
	midwest := 0
	rev := map[string]int{}
	for _, r := range rows[1:] {
		if r[col["division"]] == "Midwest" {
			midwest++
		}
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
	json.NewEncoder(os.Stdout).Encode(map[string]check{
		"midwest-regions":     {Value: strconv.Itoa(midwest)},
		"revenue-by-division": {Table: table},
	})
}
