package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// For any operation on villagers with grumpiness x <= y:
// cost = y, and afterwards the smaller grumpiness becomes 0.
// To make the whole village connected we need n-1 edges. If a component already
// contains a villager with grumpiness 0, it can be joined to another component
// that also has a 0-grumpiness villager for free. Hence we only need to pay for
// edges that create at least one zero-grumpiness villager inside each component.
//
// The cheapest way to create such a zero inside a component of size >= 2 is to
// connect its two villagers with the smallest grumpiness values; the cost of
// that edge is the larger of the two. For a sorted list, that cost is the
// second element of the pair. Doing this independently for disjoint pairs of
// villagers minimises the total cost, because the smallest villager can never
// contribute less than the next one to the total (pairing it with any larger
// villager costs at least that larger value).
//
// Thus, after sorting in non-decreasing order, pairing adjacent villagers
// minimises the sum of edge costs that create zeros. If n is even, every
// villager can be placed in such a pair and all remaining edges can be added
// for free between zero-grumpiness villagers. The total minimal cost is the
// sum of the larger element in each adjacent pair, i.e. elements at even
// indices (1-based).
//
// If n is odd, one villager remains unpaired. The cheapest way to connect this
// leftover is to attach it to any component containing a zero-grumpiness
// villager, costing exactly its own grumpiness. In the sorted order this
// leftover is the largest element (last in the array). Therefore the minimal
// total cost is the sum of elements at even indices (1-based) plus the last
// element when n is odd.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		g := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &g[i])
		}
		sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })

		var ans int64
		for i := 1; i < n; i += 2 { // add second element of each adjacent pair
			ans += g[i]
		}
		if n%2 == 1 { // leftover largest element
			ans += g[n-1]
		}
		fmt.Fprintln(out, ans)
	}
}

