package main

import (
	"bufio"
	"fmt"
	"os"
)

func dominates(a, b []int) bool {
	for i := range a {
		if a[i] <= b[i] {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	nondom := make([][]int, 0)
	for i := 0; i < n; i++ {
		p := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &p[j])
		}

		// check if any existing nondominated player dominates p
		dominated := false
		for _, q := range nondom {
			if dominates(q, p) {
				dominated = true
				break
			}
		}
		if !dominated {
			// remove players dominated by p
			newSet := nondom[:0]
			for _, q := range nondom {
				if !dominates(p, q) {
					newSet = append(newSet, q)
				}
			}
			nondom = append(newSet, p)
		}

		fmt.Fprintln(out, len(nondom))
	}
}
