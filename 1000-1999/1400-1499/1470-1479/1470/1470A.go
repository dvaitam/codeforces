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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		k := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &k[i])
			k[i]-- // convert to 0-based
		}
		c := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &c[i])
		}

		sort.Slice(k, func(i, j int) bool { return k[i] > k[j] })

		var cost int64
		ptr := 0
		for _, idx := range k {
			if ptr < m && ptr <= idx {
				cost += int64(c[ptr])
				ptr++
			} else {
				cost += int64(c[idx])
			}
		}
		fmt.Fprintln(writer, cost)
	}
}
