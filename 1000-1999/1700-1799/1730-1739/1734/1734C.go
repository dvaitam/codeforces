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
		cost := make([]int, n+1)
		for i := 1; i <= n; i++ {
			for j := i; j <= n; j += i {
				if s[j-1] == '1' {
					break
				}
				if cost[j] == 0 {
					cost[j] = i
				}
			}
		}
		ans := 0
		for i := 1; i <= n; i++ {
			if s[i-1] == '0' {
				ans += cost[i]
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
