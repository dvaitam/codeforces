package main

import (
	"bufio"
	"fmt"
	"os"
)

func canTransform(p, s string) bool {
	i, j := 0, 0
	n, m := len(p), len(s)

	for i < n {
		c := p[i]
		cntP := 0
		for i < n && p[i] == c {
			cntP++
			i++
		}

		if j >= m || s[j] != c {
			return false
		}
		cntS := 0
		for j < m && s[j] == c {
			cntS++
			j++
		}

		if cntS < cntP || cntS > 2*cntP {
			return false
		}
	}
	return j == m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var p, s string
		fmt.Fscan(in, &p)
		fmt.Fscan(in, &s)

		if canTransform(p, s) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
