package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution to problemB.txt for 1646B.
// Checks if there exist two disjoint sets of numbers painted Red and Blue
// where Sum(Red) > Sum(Blue) and Count(Red) < Count(Blue).
// It is optimal to take the smallest numbers for Blue and the largest
// numbers for Red. We iterate over possible counts of Red and compare sums.
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
		sort.Ints(a)
		// prefix sums to get sums of ranges quickly
		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + int64(a[i])
		}
		total := pref[n]
		possible := false
		// number of red elements r must satisfy 2*r+1 <= n
		for r := 1; 2*r+1 <= n; r++ {
			blueSum := pref[r+1]
			redSum := total - pref[n-r]
			if redSum > blueSum {
				possible = true
				break
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
