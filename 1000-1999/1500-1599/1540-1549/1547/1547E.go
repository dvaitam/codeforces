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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	const INF int64 = 1e18
	for ; q > 0; q-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		pos := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &pos[i])
		}
		val := make([]int64, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &val[i])
		}
		ans := make([]int64, n)
		for i := 0; i < n; i++ {
			ans[i] = INF
		}
		for i := 0; i < k; i++ {
			idx := pos[i] - 1
			if val[i] < ans[idx] {
				ans[idx] = val[i]
			}
		}
		for i := 1; i < n; i++ {
			if ans[i-1]+1 < ans[i] {
				ans[i] = ans[i-1] + 1
			}
		}
		for i := n - 2; i >= 0; i-- {
			if ans[i+1]+1 < ans[i] {
				ans[i] = ans[i+1] + 1
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, ans[i])
		}
		writer.WriteByte('\n')
	}
}
