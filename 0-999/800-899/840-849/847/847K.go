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

	var n, a, b, k, f int
	if _, err := fmt.Fscan(in, &n, &a, &b, &k, &f); err != nil {
		return
	}

	routeCost := make(map[string]int)
	total := 0
	prevTo := ""

	for i := 0; i < n; i++ {
		var s, t string
		fmt.Fscan(in, &s, &t)
		cost := a
		if i > 0 && s == prevTo {
			cost = b
		}
		total += cost
		key := s + "#" + t
		if s > t {
			key = t + "#" + s
		}
		routeCost[key] += cost
		prevTo = t
	}

	costs := make([]int, 0, len(routeCost))
	for _, v := range routeCost {
		costs = append(costs, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(costs)))

	for i := 0; i < len(costs) && i < k; i++ {
		if costs[i] > f {
			total -= costs[i] - f
		} else {
			break
		}
	}

	fmt.Fprintln(out, total)
}
