package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program attempts to solve the problem described in problemA.txt.
// It tries to construct an order of placing the given distinct weights on
// two sides of a balance so that after each placement the heavier side
// matches the sequence specified in the string S. The approach uses a
// greedy strategy with sorted weights.

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	var s string
	fmt.Fscan(in, &s)
	sort.Ints(a)

	left, right := 0, 0
	lIdx, rIdx := 0, n-1
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	prev := byte(' ')
	for i := 0; i < n; i++ {
		ch := s[i]
		if i == 0 || ch != prev {
			// take from the end when side should change
			w := a[rIdx]
			rIdx--
			if ch == 'L' {
				left += w
				fmt.Fprintln(out, w, "L")
			} else {
				right += w
				fmt.Fprintln(out, w, "R")
			}
		} else {
			// use smallest remaining weight to keep diff small
			w := a[lIdx]
			lIdx++
			if ch == 'L' {
				if left > right && w < left-right {
					right += w
					fmt.Fprintln(out, w, "R")
				} else {
					left += w
					fmt.Fprintln(out, w, "L")
				}
			} else { // ch == 'R'
				if right > left && w < right-left {
					left += w
					fmt.Fprintln(out, w, "L")
				} else {
					right += w
					fmt.Fprintln(out, w, "R")
				}
			}
		}
		prev = ch
		if (left > right) != (ch == 'L') {
			fmt.Println(-1)
			return
		}
	}
}
