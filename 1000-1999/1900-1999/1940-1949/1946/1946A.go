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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		m := (n - 1) / 2
		target := a[m] + 1
		var ops int64
		for i := m; i < n; i++ {
			if a[i] < target {
				ops += target - a[i]
			}
		}
		fmt.Fprintln(out, ops)
	}
}
