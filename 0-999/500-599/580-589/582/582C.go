package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	groups := make(map[int][]int)
	for s := 1; s < n; s++ {
		g := gcd(s, n)
		groups[g] = append(groups[g], s)
	}

	var ans int64
	for g, sList := range groups {
		maxv := make([]int, g)
		for i := 0; i < n; i++ {
			r := i % g
			if a[i] > maxv[r] {
				maxv[r] = a[i]
			}
		}
		good := make([]bool, n)
		for i := 0; i < n; i++ {
			if a[i] == maxv[i%g] {
				good[i] = true
			}
		}
		dp := make([]int, 2*n+1)
		for i := 2*n - 1; i >= 0; i-- {
			if good[i%n] {
				dp[i] = dp[i+1] + 1
			} else {
				dp[i] = 0
			}
		}
		counts := make([]int, n+1)
		for i := 0; i < n; i++ {
			L := dp[i]
			if L > n {
				L = n
			}
			counts[1]++
			if L+1 <= n {
				counts[L+1]--
			}
		}
		for k := 1; k <= n; k++ {
			counts[k] += counts[k-1]
		}
		for _, s := range sList {
			ans += int64(counts[s])
		}
	}

	fmt.Fprintln(writer, ans)
}
