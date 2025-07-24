package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Group struct {
	size  int
	money int
	id    int
}

type Table struct {
	cap  int
	id   int
	used bool
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	groups := make([]Group, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &groups[i].size, &groups[i].money)
		groups[i].id = i + 1
	}
	var k int
	fmt.Fscan(in, &k)
	tables := make([]Table, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &tables[i].cap)
		tables[i].id = i + 1
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].money == groups[j].money {
			return groups[i].size < groups[j].size
		}
		return groups[i].money > groups[j].money
	})
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].cap < tables[j].cap
	})

	results := make([][2]int, 0)
	total := 0
	for _, g := range groups {
		for j := 0; j < len(tables); j++ {
			if !tables[j].used && tables[j].cap >= g.size {
				tables[j].used = true
				total += g.money
				results = append(results, [2]int{g.id, tables[j].id})
				break
			}
		}
	}

	fmt.Fprintln(out, len(results), total)
	for _, r := range results {
		fmt.Fprintln(out, r[0], r[1])
	}
}
