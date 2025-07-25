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
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var sum int64
		for i := 0; i < n; i++ {
			sum += a[i]
		}
		ans := sum
		var suffix int64
		for i := n - 1; i > 0; i-- {
			suffix += a[i]
			if suffix > 0 {
				ans += suffix
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
