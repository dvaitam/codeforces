package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// This program solves problem G from the 1950 contest. The task is to
// remove the minimal number of songs so that the remaining songs can be
// rearranged into an "exciting" playlist, meaning adjacent songs share
// the same genre or the same writer. With n <= 16 we can use dynamic
// programming over subsets to find the largest subset that admits such an
// ordering.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)

		genres := make([]string, n)
		writers := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &genres[i], &writers[i])
		}

		adj := make([]uint16, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				if genres[i] == genres[j] || writers[i] == writers[j] {
					adj[i] |= 1 << uint(j)
				}
			}
		}

		maskCount := 1 << uint(n)
		dp := make([]uint16, maskCount)
		for i := 0; i < n; i++ {
			dp[1<<uint(i)] = 1 << uint(i)
		}

		best := 1
		for mask := 1; mask < maskCount; mask++ {
			lasts := dp[mask]
			if lasts == 0 {
				continue
			}
			sz := bits.OnesCount(uint(mask))
			if sz > best {
				best = sz
			}
			for i := 0; i < n; i++ {
				if lasts&(1<<uint(i)) == 0 {
					continue
				}
				avail := adj[i] &^ uint16(mask)
				for avail != 0 {
					j := bits.TrailingZeros16(avail)
					avail &^= 1 << uint(j)
					dp[mask|1<<uint(j)] |= 1 << uint(j)
				}
			}
		}

		fmt.Fprintln(writer, n-best)
	}
}
