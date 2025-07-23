package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const maxVal = 3000000
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	freq := make([]int64, maxVal+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
	}
	var m int
	fmt.Fscan(in, &m)
	qs := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &qs[i])
	}

	pairExact := make([]int64, maxVal+1)
	for x := 1; x <= maxVal; x++ {
		fx := freq[x]
		if fx == 0 {
			continue
		}
		for v := x; v <= maxVal; v += x {
			y := v / x
			fy := freq[y]
			if fy == 0 {
				continue
			}
			pairExact[v] += fx * fy
		}
	}

	for x := 1; x*x <= maxVal; x++ {
		if freq[x] > 0 {
			pairExact[x*x] -= freq[x]
		}
	}

	prefix := make([]int64, maxVal+1)
	for i := 1; i <= maxVal; i++ {
		prefix[i] = prefix[i-1] + pairExact[i]
	}

	totalPairs := int64(n) * int64(n-1)

	out := bufio.NewWriter(os.Stdout)
	for _, p := range qs {
		if p <= 1 {
			fmt.Fprintln(out, totalPairs)
		} else if p > maxVal {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, totalPairs-prefix[p-1])
		}
	}
	out.Flush()
}
