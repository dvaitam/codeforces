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
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		start := n - 2*k
		if start < 1 {
			start = 1
		}
		ans := int64(-1 << 63)
		for i := start; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				val := int64(i*j) - int64(k)*int64(a[i]|a[j])
				if val > ans {
					ans = val
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
