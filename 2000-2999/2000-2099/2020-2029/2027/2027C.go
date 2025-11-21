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

		ops := make(map[int64][]int64)
		for i := 2; i <= n; i++ {
			target := a[i-1] + int64(i) - 1
			ops[target] = append(ops[target], int64(i-1))
		}

		keys := make([]int64, 0, len(ops))
		for k := range ops {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })

		best := make(map[int64]int64, len(ops))
		for _, target := range keys {
			bestVal := target
			for _, delta := range ops[target] {
				next := target + delta
				candidate := next
				if v, ok := best[next]; ok {
					candidate = v
				}
				if candidate > bestVal {
					bestVal = candidate
				}
			}
			best[target] = bestVal
		}

		ans := int64(n)
		if v, ok := best[int64(n)]; ok && v > ans {
			ans = v
		}
		fmt.Fprintln(out, ans)
	}
}
