package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA.txt for 1630A (And Matching).
// It prints n/2 pairs of integers from 0..n-1 whose bitwise AND
// sums to the given k, or -1 if impossible.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		if k == n-1 {
			if n == 4 {
				fmt.Fprintln(out, -1)
				continue
			}
			used := make([]bool, n)
			// Predefined pairs to achieve sum n-1
			pairs := make([][2]int, 0, n/2)
			pairs = append(pairs, [2]int{n - 1, n - 2})
			used[n-1] = true
			used[n-2] = true
			pairs = append(pairs, [2]int{1, n - 3})
			used[1] = true
			used[n-3] = true
			pairs = append(pairs, [2]int{0, 2})
			used[0] = true
			used[2] = true
			for i := 0; i < n; i++ {
				j := n - 1 - i
				if i >= j {
					break
				}
				if used[i] || used[j] {
					continue
				}
				pairs = append(pairs, [2]int{i, j})
				used[i] = true
				used[j] = true
			}
			for _, p := range pairs {
				fmt.Fprintln(out, p[0], p[1])
			}
		} else {
			used := make([]bool, n)
			pairs := make([][2]int, 0, n/2)
			if k != 0 {
				pairs = append(pairs, [2]int{k, n - 1})
				used[k] = true
				used[n-1] = true
				pairs = append(pairs, [2]int{0, n - 1 - k})
				used[0] = true
				used[n-1-k] = true
			}
			for i := 0; i < n; i++ {
				j := n - 1 - i
				if i >= j {
					break
				}
				if used[i] || used[j] {
					continue
				}
				pairs = append(pairs, [2]int{i, j})
				used[i] = true
				used[j] = true
			}
			for _, p := range pairs {
				fmt.Fprintln(out, p[0], p[1])
			}
		}
	}
}
