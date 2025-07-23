package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	type pair struct{ cost, idx int }
	var upgrades []pair
	rating := 0
	for i := 0; i < n; i++ {
		rating += a[i] / 10
		if a[i] < 100 {
			r := a[i] % 10
			if r != 0 {
				upgrades = append(upgrades, pair{10 - r, i})
			}
		}
	}

	sort.Slice(upgrades, func(i, j int) bool { return upgrades[i].cost < upgrades[j].cost })

	for _, p := range upgrades {
		if k < p.cost {
			break
		}
		k -= p.cost
		a[p.idx] += p.cost
		rating++
	}

	total := 0
	for i := 0; i < n; i++ {
		if a[i] < 100 {
			total += (100 - a[i]) / 10
		}
	}

	if add := k / 10; add < total {
		rating += add
	} else {
		rating += total
	}

	fmt.Fprintln(out, rating)
}
