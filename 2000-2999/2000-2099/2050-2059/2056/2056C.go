package main

import (
	"bufio"
	"fmt"
	"os"
)

// Construction:
// - The value 1 appears exactly three times: at positions 1, n-1, n.
// - Every other position i (2..n-2) gets a distinct value i (so values 2..n-2).
// Only the value 1 repeats, so the longest palindromic subsequence has length 3.
// Its middle element can be any position between the first and one of the last two 1's,
// yielding (n-3)+(n-2)=2n-5 such subsequences, which is > n for all n >= 6.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		res := make([]int, n)
		for i := 0; i < n; i++ {
			if i == 0 || i == n-2 || i == n-1 {
				res[i] = 1
			} else {
				res[i] = i + 1 // distinct values, all <= n
			}
		}

		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
