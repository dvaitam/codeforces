package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
		}

		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + freq[i]
		}

		ans := 1
		// For a candidate gcd d, only numbers < 4d that are not multiples of d need erasing.
		for d := n; d >= 1; d-- {
			limit := 4*d - 1
			if limit > n {
				limit = n
			}

			less := pref[limit]
			div := 0
			if d <= limit {
				div += freq[d]
			}
			if 2*d <= limit {
				div += freq[2*d]
			}
			if 3*d <= limit {
				div += freq[3*d]
			}

			if less-div <= k {
				ans = d
				break
			}
		}

		fmt.Fprintln(out, ans)
	}
}
