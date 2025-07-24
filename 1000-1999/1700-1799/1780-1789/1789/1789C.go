package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		maxVal := n + m + 5
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		cnt := make([]int64, maxVal)
		start := make([]int, maxVal)
		for i := range start {
			start[i] = -1
		}
		for _, v := range a {
			start[v] = 0
		}
		for i := 1; i <= m; i++ {
			var p, v int
			fmt.Fscan(reader, &p, &v)
			p--
			old := a[p]
			cnt[old] += int64(i - start[old])
			start[old] = -1
			a[p] = v
			start[v] = i
		}
		for val := 1; val < maxVal; val++ {
			if start[val] != -1 {
				cnt[val] += int64(m + 1 - start[val])
			}
		}
		totalPairs := int64(m+1) * int64(m) / 2
		res := int64(2*n) * totalPairs
		var inter int64
		for val := 1; val < maxVal; val++ {
			if cnt[val] > 1 {
				inter += cnt[val] * (cnt[val] - 1) / 2
			}
		}
		res -= inter
		fmt.Fprintln(writer, res)
	}
}
