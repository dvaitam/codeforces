package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to Codeforces problem 1554A - Cherry.
// For each test case we are given n positive integers.
// The value max(a_l..a_r) * min(a_l..a_r) is maximized
// by choosing a subarray of length two, because extending
// the segment can only decrease the minimum or increase
// the maximum without improving their product.
// Thus the answer is simply the maximum product of two
// consecutive elements.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		maxProd := 0
		for i := 0; i < n-1; i++ {
			prod := a[i] * a[i+1]
			if prod > maxProd {
				maxProd = prod
			}
		}
		fmt.Fprintln(writer, maxProd)
	}
}
