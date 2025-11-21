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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		freq := make(map[int64]int64)
		var totalSum, adjustedSum int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			totalSum += x
			if x%2 == 0 {
				adjustedSum += x
			} else {
				adjustedSum += x - 1
				freq[x]++
			}
		}

		alice := adjustedSum / 2

		counts := make([]int64, 0, len(freq))
		for _, c := range freq {
			counts = append(counts, c)
		}
		sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })
		for i := 0; i < len(counts); i += 2 {
			alice += counts[i]
		}

		fmt.Fprintln(out, alice, totalSum-alice)
	}
}

