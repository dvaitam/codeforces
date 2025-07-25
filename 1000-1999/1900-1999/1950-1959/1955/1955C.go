package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution to problemC.txt (Kraken attacks the fleet).
// We alternate attacks on the first and last remaining ships.
// After k attacks we need to count how many ships have sunk.
// Instead of simulating each attack, we distribute the total number
// of attacks on the left and right ends: the left side gets ceil(k/2)
// hits and the right side gets floor(k/2) hits as long as at least two
// ships remain. Using prefix sums we determine how many ships are fully
// destroyed from each side. If only one ship remains after that,
// the remaining hits from both sides combine to possibly sink it.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
		}
		total := prefix[n]
		if k >= total {
			fmt.Fprintln(out, n)
			continue
		}

		leftHits := (k + 1) / 2
		rightHits := k / 2

		leftCount := sort.Search(n+1, func(i int) bool { return prefix[i] > leftHits }) - 1
		remLeft := leftHits - prefix[leftCount]

		j := sort.Search(n+1, func(i int) bool { return prefix[i] >= total-rightHits })
		rightCount := n - j
		remRight := rightHits - (total - prefix[j])

		sunk := leftCount + rightCount
		if sunk >= n {
			fmt.Fprintln(out, n)
			continue
		}
		if sunk == n-1 {
			midDur := a[leftCount]
			if remLeft+remRight >= midDur {
				sunk++
			}
		}
		fmt.Fprintln(out, sunk)
	}
}
