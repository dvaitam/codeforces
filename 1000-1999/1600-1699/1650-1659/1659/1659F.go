package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSorted(p []int) bool {
	for i, v := range p {
		if v != i+1 {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		for i := 0; i < n-1; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			// edges are ignored in this placeholder implementation
			_ = a
			_ = b
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if isSorted(p) {
			fmt.Fprintln(out, "Alice")
		} else {
			fmt.Fprintln(out, "Bob")
		}
	}
}
