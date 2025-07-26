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
		a := make([]int, n+1)
		prefix := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			prefix[i] = prefix[i-1] ^ a[i]
		}
		var s string
		fmt.Fscan(in, &s)
		var q int
		fmt.Fscan(in, &q)

		xor0, xor1 := 0, 0
		for i, ch := range s {
			if ch == '0' {
				xor0 ^= a[i+1]
			} else {
				xor1 ^= a[i+1]
			}
		}

		results := make([]int, 0)
		for ; q > 0; q-- {
			var tp int
			fmt.Fscan(in, &tp)
			if tp == 1 {
				var l, r int
				fmt.Fscan(in, &l, &r)
				seg := prefix[r] ^ prefix[l-1]
				xor0 ^= seg
				xor1 ^= seg
			} else {
				var g int
				fmt.Fscan(in, &g)
				if g == 0 {
					results = append(results, xor0)
				} else {
					results = append(results, xor1)
				}
			}
		}
		for i, v := range results {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
