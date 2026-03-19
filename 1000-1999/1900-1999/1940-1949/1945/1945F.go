package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		v := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &v[i])
		}
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}

		// For each k (number of mushrooms picked), positions p[1]..p[k-1] become 0.
		// We want to pick k mushrooms to maximize:
		//   (# picked with nonzero power) * min(nonzero powers among picked)
		// Nonzero means the mushroom index is NOT in {p[1]..p[k-1]}.
		// So usable = picked mushrooms whose index is in {p[k]..p[n]}.
		// We want to maximize |usable| * min(v[i] for i in usable), minimizing |usable| on ties.

		bestStrength := int64(0)
		bestCount := 0

		for k := 1; k <= n; k++ {
			// Zeroed indices: p[1]..p[k-1]
			zeroed := make(map[int]bool, k-1)
			for i := 1; i < k; i++ {
				zeroed[p[i]] = true
			}

			// Collect values of all mushrooms NOT zeroed
			var vals []int
			for i := 1; i <= n; i++ {
				if !zeroed[i] {
					vals = append(vals, v[i])
				}
			}

			// We can pick up to k mushrooms. Usable ones are those not zeroed.
			// To maximize strength, pick the top-valued non-zeroed mushrooms.
			// We can pick at most k total, but we want to maximize usable count * min.
			// Since picking zeroed mushrooms doesn't help (they contribute 0),
			// we should pick only non-zeroed mushrooms (up to k of them).
			sort.Sort(sort.Reverse(sort.IntSlice(vals)))

			maxPick := k
			if maxPick > len(vals) {
				maxPick = len(vals)
			}

			for h := 1; h <= maxPick; h++ {
				minVal := math.MaxInt64
				for j := 0; j < h; j++ {
					if vals[j] < minVal {
						minVal = vals[j]
					}
				}
				strength := int64(h) * int64(minVal)
				if strength > bestStrength || (strength == bestStrength && h < bestCount) {
					bestStrength = strength
					bestCount = h
				}
			}
		}

		fmt.Fprintln(out, bestStrength, bestCount)
	}
}
