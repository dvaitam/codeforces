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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k, x, y int
		fmt.Fscan(reader, &n, &k, &x, &y)
		a := make([]int, n)
		total := 0
		countAbove := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			total += a[i]
			if a[i] > k {
				countAbove++
			}
		}

		// Option without even distribution
		cost := countAbove * x

		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		prefix := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
		}

		// Find minimal number of removals to make total <= k*n
		target := k * n
		r := sort.Search(n+1, func(r int) bool {
			remaining := total - prefix[r]
			return remaining <= target
		})
		if r <= n {
			costDist := r*x + y
			if costDist < cost {
				cost = costDist
			}
		}

		fmt.Fprintln(writer, cost)
	}
}
