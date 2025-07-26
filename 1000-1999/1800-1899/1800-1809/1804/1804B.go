package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt from contest 1804.
// Greedy approach: process patients in order, opening a new pack whenever
// the current one has no doses left or expires before the next patient can
// be vaccinated. When opening a pack for patient i at time t[i]+w, its doses
// remain usable until t[i]+w+d. Each patient can wait up to w moments, so
// this schedule is always feasible. We keep using the current pack while the
// next patient's arrival time does not exceed the expiry and we still have
// doses left.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k, d, w int
		fmt.Fscan(in, &n, &k, &d, &w)
		t := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &t[i])
		}

		packs := 0
		expire := -1
		left := 0
		for i := 0; i < n; i++ {
			if t[i] > expire || left == 0 {
				packs++
				expire = t[i] + w + d
				left = k
			}
			left--
		}
		fmt.Fprintln(out, packs)
	}
}
