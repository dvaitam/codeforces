package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// Given an array, we need a non-empty subsequence with the
// same contrast (sum of absolute differences of consecutive
// elements) as the whole array but with minimal length.
//
// Removing an element a[i] does not change the contrast if
// it lies between its neighbours in value, i.e. when
// |prev-a[i]| + |a[i]-next| = |prev-next|. Therefore the
// minimal subsequence keeps only the endpoints and the points
// where the direction of monotonicity changes (local maxima
// and minima). After building this compressed sequence, if it
// contains exactly two equal numbers we can reduce it to one
// element since the contrast is zero.
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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		b := make([]int, 0, n)
		b = append(b, a[0])
		for i := 1; i < n-1; i++ {
			prev := b[len(b)-1]
			curr := a[i]
			next := a[i+1]
			if (prev <= curr && curr <= next) || (prev >= curr && curr >= next) {
				continue
			}
			b = append(b, curr)
		}
		b = append(b, a[n-1])
		ans := len(b)
		if ans == 2 && b[0] == b[1] {
			ans = 1
		}
		fmt.Fprintln(out, ans)
	}
}
