package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		cnt := make([]int, n+1)
		last := make([]int, n+1)
		for i := range last {
			last[i] = -1
		}
		for i := 1; i <= n; i++ {
			var c int
			fmt.Fscan(in, &c)
			p := i % 2
			if last[c] == -1 || last[c] != p {
				cnt[c]++
				last[c] = p
			}
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, cnt[i])
		}
		fmt.Fprintln(out)
	}
}
