package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func can(a []int, k int, d int) bool {
	n := len(a)
	dp := make([]bool, n+1)
	pref := make([]int, n+1)
	dp[0] = true
	pref[0] = 1
	j := 0
	for i := 1; i <= n; i++ {
		if i >= k {
			for j < i && a[i-1]-a[j] > d {
				j++
			}
			if j <= i-k {
				sum := pref[i-k]
				if j > 0 {
					sum -= pref[j-1]
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
	return dp[n]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	sort.Ints(a)
	if k <= 1 {
		fmt.Fprintln(writer, 0)
		return
	}

	low, high := 0, a[n-1]-a[0]
	for low < high {
		mid := (low + high) / 2
		if can(a, k, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}

	fmt.Fprintln(writer, low)
}
