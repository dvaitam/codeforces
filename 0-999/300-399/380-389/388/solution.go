package main

import (
	"fmt"
	"sort"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	var cielSum, jiroSum int
	var middles []int

	for i := 0; i < n; i++ {
		var s int
		fmt.Scan(&s)
		cards := make([]int, s)
		for j := 0; j < s; j++ {
			fmt.Scan(&cards[j])
		}

		// In optimal play for a pile of size s:
		// Ciel is guaranteed the top s/2 cards.
		// Jiro is guaranteed the bottom s/2 cards.
		// If s is odd, the middle card is contested.

		count := s / 2
		for j := 0; j < count; j++ {
			cielSum += cards[j]
			jiroSum += cards[s-1-j]
		}

		if s%2 == 1 {
			middles = append(middles, cards[s/2])
		}
	}

	// The middle cards from all odd-sized piles form a pool of available cards.
	// Since Ciel goes first and the parity of moves to expose a middle card is consistent,
	// players will take turns picking the largest available middle card, starting with Ciel.
	sort.Sort(sort.Reverse(sort.IntSlice(middles)))

	for i, val := range middles {
		if i%2 == 0 {
			cielSum += val
		} else {
			jiroSum += val
		}
	}

	fmt.Println(cielSum, jiroSum)
}
