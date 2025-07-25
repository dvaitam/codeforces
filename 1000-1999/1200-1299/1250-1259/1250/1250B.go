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
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	teams := make([]int, k)
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		teams[t-1]++
	}
	sort.Ints(teams)
	best := int64(^uint64(0) >> 1)
	maxPair := 0
	for pairs := 0; pairs <= k/2; pairs++ {
		if pairs > 0 {
			sum := teams[pairs-1] + teams[k-pairs]
			if sum > maxPair {
				maxPair = sum
			}
		}
		rides := k - pairs
		largest := 0
		if pairs <= k-pairs-1 {
			largest = teams[k-pairs-1]
		}
		cap := largest
		if maxPair > cap {
			cap = maxPair
		}
		cost := int64(rides) * int64(cap)
		if cost < best {
			best = cost
		}
	}
	fmt.Println(best)
}
