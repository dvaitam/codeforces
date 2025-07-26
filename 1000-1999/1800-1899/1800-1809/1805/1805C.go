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
		ks := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &ks[i])
		}
		sort.Ints(ks)
		for i := 0; i < m; i++ {
			var a, b, c int64
			fmt.Fscan(reader, &a, &b, &c)
			d := 4 * a * c
			idx := sort.SearchInts(ks, int(b))
			found := false
			ans := 0
			if idx < len(ks) {
				diff := int64(ks[idx]) - b
				if diff*diff < d {
					found = true
					ans = ks[idx]
				}
			}
			if !found && idx > 0 {
				diff := int64(ks[idx-1]) - b
				if diff*diff < d {
					found = true
					ans = ks[idx-1]
				}
			}
			if found {
				fmt.Fprintln(writer, "YES")
				fmt.Fprintln(writer, ans)
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
