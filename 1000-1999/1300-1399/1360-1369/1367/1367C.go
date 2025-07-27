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
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)

		right := make([]int, n)
		next := n * 2
		for i := n - 1; i >= 0; i-- {
			if s[i] == '1' {
				next = i
			}
			right[i] = next
		}

		last := -n * 2
		ans := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				last = i
				continue
			}
			if i-last > k && right[i]-i > k {
				ans++
				last = i
			}
		}
		fmt.Fprintln(out, ans)
	}
}
