package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		// Compress values to ranks 1..n
		sortedVals := append([]int(nil), arr...)
		sort.Ints(sortedVals)
		rank := make(map[int]int, n)
		for i, v := range sortedVals {
			rank[v] = i + 1
		}
		p := make([]int, n)
		for i, v := range arr {
			p[i] = rank[v]
		}

		// Precompute total length of all subarrays
		var totalLength int64
		for k := 1; k <= n; k++ {
			totalLength += int64(k) * int64(n-k+1)
		}

		// Each subarray contributes at least one boundary at its end
		totalBoundaries := int64(n * (n + 1) / 2)

		rightMin := make([]int, n)
		leftVals := make([]int, n)
		INF := n + 5
		for i := 1; i < n; i++ {
			m := INF
			for r := i; r < n; r++ {
				if p[r] < m {
					m = p[r]
				}
				rightMin[r-i] = m
			}
			m = 0
			for l := i - 1; l >= 0; l-- {
				if p[l] > m {
					m = p[l]
				}
				leftVals[l] = m
			}
			rPtr := 0
			limit := n - i
			for l := 0; l < i; l++ {
				left := leftVals[l]
				for rPtr < limit && rightMin[rPtr] > left {
					rPtr++
				}
				totalBoundaries += int64(rPtr)
			}
		}

		result := totalLength - totalBoundaries
		fmt.Fprintln(writer, result)
	}
}
