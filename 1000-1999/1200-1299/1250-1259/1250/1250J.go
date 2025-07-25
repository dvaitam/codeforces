package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxRows(counts []int64, L int64) int64 {
	if L == 0 {
		return 0
	}
	var rows, rem int64
	for _, x := range counts {
		if rem > 0 {
			need := L - rem
			if x >= need {
				rows++
				x -= need
				rem = 0
			} else {
				rem = 0
			}
		}
		rows += x / L
		rem = x % L
	}
	return rows
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		counts := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &counts[i])
			sum += counts[i]
		}
		if sum < k {
			fmt.Fprintln(writer, 0)
			continue
		}
		lo := int64(0)
		hi := sum/k + 1
		for hi-lo > 1 {
			mid := (lo + hi) / 2
			if maxRows(counts, mid) >= k {
				lo = mid
			} else {
				hi = mid
			}
		}
		fmt.Fprintln(writer, lo*k)
	}
}
