package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	pairs := make([]pair, n)
	prefixMin := make([]int, n)
	suffixMax := make([]int, n)

	for ; q > 0; q-- {
		var pos, x int
		fmt.Fscan(in, &pos, &x)
		pos--
		a[pos] = x

		for i := 0; i < n; i++ {
			pairs[i] = pair{val: a[i], idx: i}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })

		prefixMin[0] = pairs[0].idx
		for i := 1; i < n; i++ {
			if pairs[i].idx < prefixMin[i-1] {
				prefixMin[i] = pairs[i].idx
			} else {
				prefixMin[i] = prefixMin[i-1]
			}
		}
		suffixMax[n-1] = pairs[n-1].idx
		for i := n - 2; i >= 0; i-- {
			if pairs[i].idx > suffixMax[i+1] {
				suffixMax[i] = pairs[i].idx
			} else {
				suffixMax[i] = suffixMax[i+1]
			}
		}
		weight := 1
		for i := 0; i < n-1; i++ {
			if prefixMin[i] > suffixMax[i+1] {
				weight++
			}
		}
		fmt.Fprintln(out, weight)
	}
}
