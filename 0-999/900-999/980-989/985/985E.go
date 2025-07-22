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

	var n, k int
	var d int
	if _, err := fmt.Fscan(in, &n, &k, &d); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)

	dp := make([]bool, n+1)
	pref := make([]int, n+1)
	dp[0] = true
	pref[0] = 1

	l := 0
	for i := 1; i <= n; i++ {
		for l < i && arr[i-1]-arr[l] > d {
			l++
		}
		if i >= k {
			left := l
			right := i - k
			if left <= right {
				sum := pref[right]
				if left > 0 {
					sum -= pref[left-1]
				}
				if sum > 0 {
					dp[i] = true
				}
			}
		}
		pref[i] = pref[i-1]
		if dp[i] {
			pref[i]++
		}
	}

	if dp[n] {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
