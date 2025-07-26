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
		s := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}
		f := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &f[i])
		}
		prev := 0
		for i := 0; i < n; i++ {
			start := s[i]
			if prev > start {
				start = prev
			}
			duration := f[i] - start
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, duration)
			prev = f[i]
		}
		if t > 0 {
			fmt.Fprintln(out)
		}
	}
}
