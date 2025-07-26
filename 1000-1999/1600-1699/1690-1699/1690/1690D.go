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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		cnt := 0
		for i := 0; i < k; i++ {
			if s[i] == 'W' {
				cnt++
			}
		}
		ans := cnt
		for i := k; i < n; i++ {
			if s[i-k] == 'W' {
				cnt--
			}
			if s[i] == 'W' {
				cnt++
			}
			if cnt < ans {
				ans = cnt
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
