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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		if s[0] != t[0] || s[n-1] != t[n-1] {
			fmt.Fprintln(out, -1)
			continue
		}

		a := make([]int, 0)
		b := make([]int, 0)
		for i := 0; i < n-1; i++ {
			if s[i] != s[i+1] {
				a = append(a, i+1)
			}
			if t[i] != t[i+1] {
				b = append(b, i+1)
			}
		}

		if len(a) != len(b) {
			fmt.Fprintln(out, -1)
			continue
		}

		res := 0
		for i := range a {
			if a[i] > b[i] {
				res += a[i] - b[i]
			} else {
				res += b[i] - a[i]
			}
		}
		fmt.Fprintln(out, res)
	}
}
