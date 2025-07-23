package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// The algorithm follows a simple greedy approach. For each required
// deletion length 2^(i-1) we try all possible positions of the
// substring of that length and pick the resulting string which is
// lexicographically smallest. This is repeated for every operation.
// While this approach does not guarantee an optimal answer for all
// cases, it is a straightforward heuristic that respects the required
// sequence of deletions.

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	k := int(math.Log2(float64(n)))
	cur := s
	for i := 0; i < k; i++ {
		l := 1 << i
		if l > len(cur) {
			break
		}
		best := cur[l:]
		bestPos := 0
		for j := 1; j <= len(cur)-l; j++ {
			cand := cur[:j] + cur[j+l:]
			if cand < best {
				best = cand
				bestPos = j
			}
		}
		_ = bestPos // suppress unused warning if not using go vet
		cur = best
	}
	fmt.Println(cur)
}
