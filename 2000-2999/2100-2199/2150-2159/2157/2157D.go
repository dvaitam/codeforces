package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var n int
		var l, r int64
		fmt.Fscan(in, &n, &l, &r)

		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
		}
		sort.Ints(a)

		pref := make([]int64, n+1)
		for j := 0; j < n; j++ {
			pref[j+1] = pref[j] + int64(a[j])
		}

		var maxAns int64 = 0

		// We iterate k = (count of +1) - (count of -1)
		// k ranges from -n to n
		for k := -n; k <= n; k++ {
			// We want to use as many elements as possible.
			// p + m = u
			// p - m = k
			// => 2p = u + k => u + k must be even.
			// We try u = n. If n+k is odd, we must use u = n-1.
			
			u := n
			if (n+k)%2 != 0 {
				u = n - 1
			}
			
			// Check if valid (u >= |k| basically)
			if u < 0 {
				continue
			}
			
			p := (u + k) / 2
			m := (u - k) / 2
			
			if p < 0 || m < 0 {
				continue
			}

			// Calculate term: sum(largest p) - sum(smallest m)
			// Largest p: a[n-p]...a[n-1]
			sumP := pref[n] - pref[n-p]
			// Smallest m: a[0]...a[m-1]
			sumM := pref[m]
			
			term := sumP - sumM
			
			var currentAns int64
			if k >= 0 {
				currentAns = term - int64(k)*r
			} else {
				currentAns = term - int64(k)*l
			}
			
			if currentAns > maxAns {
				maxAns = currentAns
			}
		}

		fmt.Fprintln(out, maxAns)
	}
}
