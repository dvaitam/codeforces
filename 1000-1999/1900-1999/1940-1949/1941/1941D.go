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
		var n, m, x int
		fmt.Fscan(in, &n, &m, &x)
		x--
		cur := make([]bool, n)
		cur[x] = true
		for i := 0; i < m; i++ {
			var r int
			var c string
			fmt.Fscan(in, &r, &c)
			next := make([]bool, n)
			for pos := 0; pos < n; pos++ {
				if !cur[pos] {
					continue
				}
				if c == "0" || c == "?" {
					to := (pos + r) % n
					next[to] = true
				}
				if c == "1" || c == "?" {
					to := pos - r
					to %= n
					if to < 0 {
						to += n
					}
					next[to] = true
				}
			}
			cur = next
		}
		var res []int
		for i := 0; i < n; i++ {
			if cur[i] {
				res = append(res, i+1)
			}
		}
		fmt.Fprintln(out, len(res))
		for i, v := range res {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
