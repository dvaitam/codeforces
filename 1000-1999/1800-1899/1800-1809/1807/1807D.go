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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		a := make([]int64, n+1)
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			prefix[i] = prefix[i-1] + a[i]
		}
		total := prefix[n]
		for ; q > 0; q-- {
			var l, r int
			var k int64
			fmt.Fscan(reader, &l, &r, &k)
			newSum := total - (prefix[r] - prefix[l-1]) + int64(r-l+1)*k
			if newSum%2 != 0 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
