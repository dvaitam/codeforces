package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// Precompute next occurrence of each position
	next := make([]int, n)
	last := make(map[int]int)
	for i := n - 1; i >= 0; i-- {
		if v, ok := last[a[i]]; ok {
			next[i] = v
		} else {
			next[i] = n
		}
		last[a[i]] = i
	}

	hand := make([]int, k)
	for i := 0; i < k; i++ {
		hand[i] = i
	}
	ptr := k
	coins := 0

	for len(hand) > 0 {
		// build map type->indices in hand
		groups := make(map[int][]int)
		for _, idx := range hand {
			groups[a[idx]] = append(groups[a[idx]], idx)
		}
		var choose [2]int
		chosen := false
		for _, idxs := range groups {
			if len(idxs) >= 2 {
				// choose two indices with farthest next occurrence
				sort.Slice(idxs, func(i, j int) bool { return next[idxs[i]] > next[idxs[j]] })
				choose[0], choose[1] = idxs[0], idxs[1]
				chosen = true
				break
			}
		}
		if !chosen {
			sort.Slice(hand, func(i, j int) bool { return next[hand[i]] > next[hand[j]] })
			choose[0], choose[1] = hand[0], hand[1]
		}
		// remove chosen indices from hand
		newHand := make([]int, 0, len(hand)-2)
		removed := 0
		for _, idx := range hand {
			if idx == choose[0] || idx == choose[1] {
				removed++
				if removed == 2 {
					// skip both
					continue
				}
				continue
			}
			newHand = append(newHand, idx)
		}
		hand = newHand
		if a[choose[0]] == a[choose[1]] {
			coins++
		}
		// draw next two
		if ptr < n {
			hand = append(hand, ptr)
			ptr++
		}
		if ptr < n {
			hand = append(hand, ptr)
			ptr++
		}
	}
	fmt.Println(coins)
}
