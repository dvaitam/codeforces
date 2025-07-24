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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		t1 := make([]int, n)
		t0 := make([]int, n)
		for i := 1; i < n; i++ {
			if s[i-1] == '1' {
				t1[i] = t1[i-1] + 1
			} else {
				t1[i] = 0
			}
			if s[i-1] == '0' {
				t0[i] = t0[i-1] + 1
			} else {
				t0[i] = 0
			}
		}
		for i := 1; i < n; i++ {
			x := i + 1
			ans := x - t1[i] - t0[i]
			if i > 1 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, ans)
		}
		out.WriteByte('\n')
	}
}
