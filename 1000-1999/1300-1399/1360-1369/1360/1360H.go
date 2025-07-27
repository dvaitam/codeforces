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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		removed := make([]uint64, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			var val uint64
			for j := 0; j < m; j++ {
				val <<= 1
				if s[j] == '1' {
					val |= 1
				}
			}
			removed[i] = val
		}
		sort.Slice(removed, func(i, j int) bool { return removed[i] < removed[j] })
		total := uint64(1) << uint(m)
		target := (total - uint64(n) - 1) / 2
		candidate := target
		for _, r := range removed {
			if r <= candidate {
				candidate++
			}
		}
		ans := make([]byte, m)
		for i := m - 1; i >= 0; i-- {
			if candidate&1 == 1 {
				ans[i] = '1'
			} else {
				ans[i] = '0'
			}
			candidate >>= 1
		}
		fmt.Fprintln(writer, string(ans))
	}
}
