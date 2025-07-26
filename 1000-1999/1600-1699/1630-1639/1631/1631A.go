package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for 1631A.
// For each index, we can optionally swap a[i] with b[i]. To minimize
// max(a) * max(b), we should ensure that a[i] <= b[i] by swapping when
// necessary. Then the answer is simply the product of the maximum values
// in arrays a and b.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		for i := 0; i < n; i++ {
			if a[i] > b[i] {
				a[i], b[i] = b[i], a[i]
			}
		}
		maxA, maxB := 0, 0
		for i := 0; i < n; i++ {
			if a[i] > maxA {
				maxA = a[i]
			}
			if b[i] > maxB {
				maxB = b[i]
			}
		}
		fmt.Fprintln(out, maxA*maxB)
	}
}
