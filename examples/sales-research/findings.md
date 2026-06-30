---
title: Findings
summary: Findings about the 2025 regional sales data, each backed by a re-runnable check.
---
# Findings

## Regional structure

{{< claim value="4" check="midwest-regions" source="data/sales.csv" asof="2026-06-30" >}}
The **Midwest** division is served by **4** distinct sales regions — the most of any division.
```sh
# evidence/midwest-regions/run
awk -F, 'NR>1 && $2=="Midwest"{n++} END{print n+0}' data/sales.csv
```
```result
4
```
{{< /claim >}}

{{< claim check="revenue-by-division" source="data/sales.csv" asof="2026-06-30" >}}
The **Midwest** also leads 2025 revenue at **11,400,000** — ahead of the Northeast's 9,000,000.
```sh
# evidence/revenue-by-division/run
echo "division,total"
awk -F, 'NR>1{r[$2]+=$3} END{for(d in r) print d","r[d]}' data/sales.csv | sort -t, -k2,2nr
```
```result
division,total
Midwest,11400000
Northeast,9000000
West,8000000
South,7900000
```
{{< /claim >}}

## A claim still needing support

{{< claim >}}
Revenue grew **18% year over year** across every division — *this assertion has no
evidence yet, so `check` flags it.*
{{< /claim >}}
