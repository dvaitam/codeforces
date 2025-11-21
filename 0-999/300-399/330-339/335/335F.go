package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	prices := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &prices[i])
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i] < prices[j]
	})

	var total int64
	freeIdx := n - 2
	for i := n - 1; i >= 0; i-- {
		total += prices[i]
		if freeIdx >= i {
			freeIdx = i - 1
		}
		for freeIdx >= 0 && prices[freeIdx] >= prices[i] {
			freeIdx--
		}
		if freeIdx >= 0 {
			freeIdx--
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, total)
	out.Flush()
}
