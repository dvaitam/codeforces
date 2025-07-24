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
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := 0
		best := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				best = a[i]
			} else {
				ans += a[i]
				if best > a[i] {
					ans += best - a[i]
					best = a[i]
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
