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
		a := make([]int, n)
		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] <= n {
				freq[a[i]]++
			}
		}
		pref := make([]int, n+2)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + freq[i]
		}
		need := (n + k + 1) / 2
		bestL, bestR := 1, n
		r := 1
		for l := 1; l <= n; l++ {
			for r <= n && pref[r]-pref[l-1] < need {
				r++
			}
			if r > n {
				break
			}
			if r-l < bestR-bestL {
				bestL, bestR = l, r
			}
		}
		fmt.Fprintln(out, bestL, bestR)
		start := 1
		balance := 0
		count := 0
		for i := 1; i <= n; i++ {
			if a[i-1] >= bestL && a[i-1] <= bestR {
				balance++
			} else {
				balance--
			}
			if balance > 0 && count < k-1 {
				fmt.Fprintln(out, start, i)
				start = i + 1
				count++
				balance = 0
			}
		}
		fmt.Fprintln(out, start, n)
	}
}
