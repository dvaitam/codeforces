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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n <= 2 {
			fmt.Fprintln(out, -1)
			continue
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		var sum int64
		for _, v := range a {
			sum += v
		}
		need := n/2 + 1
		xs := make([]int64, n-1)
		twoN := int64(2 * n)
		for i := 0; i < n-1; i++ {
			val := twoN*a[i] - sum + 1
			if val < 0 {
				val = 0
			}
			xs[i] = val
		}
		sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
		if need > len(xs) {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, xs[need-1])
		}
	}
}
