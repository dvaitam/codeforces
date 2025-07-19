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
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n == 0 || a[n-1] != 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		// build operations
		ans := make([]int, 0, n)
		prev := 0
		for i := 0; i < n; i++ {
			if a[i] == 0 {
				// first, the position p = i - prev
				ans = append(ans, i-prev)
				// then zeros for each j in [prev, i)
				for j := prev; j < i; j++ {
					ans = append(ans, 0)
				}
				prev = i + 1
			}
		}
		// reverse ans
		for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
			ans[i], ans[j] = ans[j], ans[i]
		}
		fmt.Fprintln(writer, "YES")
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
