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
		var n, k int
		var q int
		fmt.Fscan(reader, &n, &k, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var ans int64
		cur := 0
		for i := 0; i < n; i++ {
			if a[i] <= q {
				cur++
			} else {
				if cur >= k {
					l := int64(cur - k + 1)
					ans += l * (l + 1) / 2
				}
				cur = 0
			}
		}
		if cur >= k {
			l := int64(cur - k + 1)
			ans += l * (l + 1) / 2
		}
		fmt.Fprintln(writer, ans)
	}
}
