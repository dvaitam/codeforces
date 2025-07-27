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
		var n, k1, k2 int
		fmt.Fscan(in, &n, &k1, &k2)
		var s string
		fmt.Fscan(in, &s)

		prev := 0
		total := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				prev = 0
			} else {
				cur := k2 - prev
				if cur > k1 {
					cur = k1
				}
				if cur < 0 {
					cur = 0
				}
				total += cur
				prev = cur
			}
		}
		fmt.Fprintln(out, total)
	}
}
