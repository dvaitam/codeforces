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
		var n int
		var k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if a[i] <= k {
				b[i] = 1
			} else {
				b[i] = -1
			}
		}

		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + b[i]
		}

		suffixMax := make([]int, n+2)
		suffixMax[n] = pref[n]
		for i := n - 1; i >= 0; i-- {
			suffixMax[i] = pref[i]
			if suffixMax[i+1] > suffixMax[i] {
				suffixMax[i] = suffixMax[i+1]
			}
		}

		minPref := make([]int, n+1)
		minPref[0] = pref[0]
		for i := 1; i <= n; i++ {
			minPref[i] = pref[i]
			if minPref[i-1] < minPref[i] {
				minPref[i] = minPref[i-1]
			}
		}

		prefGoodExists := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			prefGoodExists[i] = prefGoodExists[i-1] || (pref[i] >= 0)
		}

		ans := false

		// Case 1: first and second segments good
		for l := 1; l <= n-2 && !ans; l++ {
			if pref[l] >= 0 {
				if suffixMax[l+1] >= pref[l] {
					ans = true
				}
			}
		}

		// Case 2: second and third segments good
		if !ans {
			for r := 2; r <= n-1 && !ans; r++ {
				if pref[n] >= pref[r] && minPref[r-1] <= pref[r] {
					ans = true
				}
			}
		}

		// Case 3: first and third segments good
		if !ans {
			for j := 3; j <= n && !ans; j++ {
				suffixSum := pref[n] - pref[j-1]
				if suffixSum >= 0 && prefGoodExists[j-2] {
					ans = true
				}
			}
		}

		if ans {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
