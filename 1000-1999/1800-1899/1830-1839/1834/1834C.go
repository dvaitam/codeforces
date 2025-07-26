package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t++ {
		var n int
		fmt.Fscan(in, &n)
		var s, r string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &r)
		a := 0
		b := 0
		for i := 0; i < n; i++ {
			if s[i] != r[i] {
				a++
			}
			if s[i] != r[n-1-i] {
				b++
			}
		}
		ans1 := 2*a - a%2
		ans2 := 2*b - (1 - b%2)
		if ans2 < 0 {
			ans2 = 1
		}
		fmt.Fprintln(out, min(ans1, ans2))
	}
}
