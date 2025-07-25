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
		var n int
		var m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// indices sorted by price
		idx := make([]int, n)
		for i := range idx {
			idx[i] = i
		}
		sort.Slice(idx, func(i, j int) bool { return a[idx[i]] < a[idx[j]] })

		x := make([]int64, n)
		left := k
		for _, id := range idx {
			if left == 0 {
				break
			}
			take := m
			if left < take {
				take = left
			}
			x[id] = take
			left -= take
		}
		var prefix, ans int64
		for i := 0; i < n; i++ {
			ans += x[i] * (int64(a[i]) + prefix)
			prefix += x[i]
		}
		fmt.Fprintln(writer, ans)
	}
}
