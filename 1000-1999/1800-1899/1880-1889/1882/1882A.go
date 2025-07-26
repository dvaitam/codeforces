package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problem A from contest 1882.
// For a given array a_1..a_n, we construct the lexicographically
// smallest increasing sequence b_1..b_n of positive integers such
// that b_i != a_i for each i. Greedily pick the smallest valid
// integer at each step. The final value of b_n is then minimal.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}
		cur := 1
		for _, x := range a {
			if cur == x {
				cur++
			}
			cur++
		}
		fmt.Fprintln(writer, cur-1)
	}
}
