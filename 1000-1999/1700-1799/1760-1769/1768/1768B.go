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
		var n, k int
		fmt.Fscan(in, &n, &k)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		expected := 1
		for _, v := range p {
			if v == expected {
				expected++
			}
		}
		m := expected - 1
		remaining := n - m
		if remaining <= 0 {
			fmt.Fprintln(out, 0)
		} else {
			ans := (remaining + k - 1) / k
			fmt.Fprintln(out, ans)
		}
	}
}
