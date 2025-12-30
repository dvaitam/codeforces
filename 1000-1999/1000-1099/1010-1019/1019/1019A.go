package main

import (
	"fmt"
	"sort"
)

// Pair struct stores the cost to buy a vote and the candidate index currently holding it.
type Pair struct {
	Cost      int
	Candidate int
}

func main() {
	var n, m int
	// Check for input errors, though strictly not necessary for competitive programming
	if _, err := fmt.Scan(&n, &m); err != nil {
		return
	}

	var c []Pair
	for i := 0; i < n; i++ {
		var p, cost int
		fmt.Scan(&p, &cost)
		// We only care about votes not already cast for candidate 1 (index 1 in input)
		if p != 1 {
			// Store candidate as 0-indexed (p-1) for easier array access later
			c = append(c, Pair{Cost: cost, Candidate: p - 1})
		}
	}

	// Sort votes by Cost descending.
	// C++ uses default pair sort (ascending) + reverse.
	// Go's sort.Slice is used here to replicate that order (Cost desc, then Candidate desc).
	sort.Slice(c, func(i, j int) bool {
		if c[i].Cost != c[j].Cost {
			return c[i].Cost > c[j].Cost
		}
		return c[i].Candidate > c[j].Candidate
	})

	// Initialize answer with a large value (roughly 9e18, fits in int64)
	var ans int64 = 9000000000000000000

	// Iterate through 'i', representing the target number of votes we want Candidate 1 to achieve.
	for i := 1; i <= n; i++ {
		// Track vote counts for opponent candidates in this iteration
		cnt := make([]int, m)
		
		// 'rest' represents the maximum number of votes we can allow ALL opponents to keep combined.
		// If Candidate 1 must win with 'i' votes, opponents hold 'N - i' votes max.
		rest := n - i
		var need int64 = 0

		// Iterate through votes, starting from the most expensive
		for _, vote := range c {
			// We try to save money by letting the opponent keep the vote if:
			// 1. We still have "capacity" in the opponent pool (rest > 0)
			// 2. This specific opponent won't reach 'i' votes (strictly less than i to ensure a win)
			if rest > 0 && cnt[vote.Candidate]+1 < i {
				cnt[vote.Candidate]++
				rest--
			} else {
				// Otherwise, we are forced to buy this vote
				need += int64(vote.Cost)
			}
		}

		if need < ans {
			ans = need
		}
	}

	fmt.Println(ans)
}
