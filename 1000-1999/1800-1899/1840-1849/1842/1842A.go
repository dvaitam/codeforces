package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It compares the total ability values of both players' monsters
// and declares the winner accordingly.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		sumA := 0
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			sumA += v
		}
		sumB := 0
		for i := 0; i < m; i++ {
			var v int
			fmt.Fscan(in, &v)
			sumB += v
		}
		if sumA > sumB {
			fmt.Fprintln(out, "Tsondu")
		} else if sumA < sumB {
			fmt.Fprintln(out, "Tenzing")
		} else {
			fmt.Fprintln(out, "Draw")
		}
	}
}
