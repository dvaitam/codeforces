package main

import (
	"bufio"
	"fmt"
	"os"
)

func canCreate(n, m int) bool {
	if n == m {
		return true
	}
	if n < m || n%3 != 0 {
		return false
	}
	a := n / 3
	return canCreate(a, m) || canCreate(n-a, m)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		if canCreate(n, m) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
