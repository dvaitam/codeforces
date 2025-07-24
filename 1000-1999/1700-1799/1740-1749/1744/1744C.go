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
		var c string
		fmt.Fscan(in, &n, &c)
		var s string
		fmt.Fscan(in, &s)
		if c == "g" {
			fmt.Fprintln(out, 0)
			continue
		}
		b := []byte(s)
		tbytes := make([]byte, 2*n)
		copy(tbytes, b)
		copy(tbytes[n:], b)
		next := make([]int, 2*n)
		pos := 2 * n
		for i := 2*n - 1; i >= 0; i-- {
			if tbytes[i] == 'g' {
				pos = i
			}
			if pos == 2*n {
				next[i] = 2 * n
			} else {
				next[i] = pos - i
			}
		}
		ans := 0
		for i := 0; i < n; i++ {
			if b[i] == c[0] {
				if next[i] > ans {
					ans = next[i]
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
