package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// For each test case we are given two arrays a and b. We may
// swap a[i] and b[i] for any i. The task is to check if we can
// perform some swaps so that a[n] becomes the maximum value of
// array a and b[n] becomes the maximum value of array b.
//
// Observation: after fixing whether we swap the last pair or not,
// the desired maxima become a[n] and b[n]. For every other pair
// (a[i], b[i]) we only need to know if it can be oriented so that
// its values do not exceed these maxima. If for each i < n there
// exists an orientation with a[i] <= maxA and b[i] <= maxB, then
// such a configuration of swaps works.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		if canMakeMax(a, b, false) || canMakeMax(a, b, true) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

// canMakeMax checks if we can orient all pairs so that a[n-1] and
// b[n-1] (possibly swapped when swapLast is true) are the maxima of
// arrays a and b.
func canMakeMax(a, b []int, swapLast bool) bool {
	n := len(a)
	maxA, maxB := a[n-1], b[n-1]
	if swapLast {
		maxA, maxB = maxB, maxA
	}
	for i := 0; i < n-1; i++ {
		ai, bi := a[i], b[i]
		if !(ai <= maxA && bi <= maxB) && !(bi <= maxA && ai <= maxB) {
			return false
		}
	}
	return true
}
