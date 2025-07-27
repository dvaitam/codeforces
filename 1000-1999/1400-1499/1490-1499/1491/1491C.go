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
		s := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &s[i])
		}
		add := make([]int64, n+2)
		var ans int64
		for i := 1; i <= n; i++ {
			if add[i] < s[i]-1 {
				ans += s[i] - 1 - add[i]
				add[i] = s[i] - 1
			}
			leftover := add[i] - (s[i] - 1)
			add[i+1] += leftover
			limit := i + int(s[i])
			if limit > n {
				limit = n
			}
			for j := i + 2; j <= limit; j++ {
				add[j]++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
