package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// Count frequencies of each value
	freqMap := make(map[int]int)
	for _, v := range a {
		freqMap[v]++
	}
	counts := make([]int, 0, len(freqMap))
	for _, c := range freqMap {
		counts = append(counts, c)
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })
	d := len(counts)
	L := counts[0]

	// Compute R[j] = number of values with count >= j
	R := make([]int, L)
	for j := 1; j <= L; j++ {
		cnt := 0
		for _, c := range counts {
			if c >= j {
				cnt++
			} else {
				break
			}
		}
		R[j-1] = cnt
	}
	// Prefix sums B[i] = sum_{t=1..i} R[t]
	B := make([]int, L)
	sum := 0
	for i := 0; i < L; i++ {
		sum += R[i]
		B[i] = sum
	}

	// Precompute number of partitions with parts <= k for all k <= d and totals <= n
	part := make([][]int, d+1)
	for i := range part {
		part[i] = make([]int, n+1)
	}
	for k := 0; k <= d; k++ {
		part[k][0] = 1
	}
	for k := 1; k <= d; k++ {
		for i := 0; i <= n; i++ {
			val := part[k-1][i]
			if i >= k {
				val += part[k][i-k]
				if val >= MOD {
					val -= MOD
				}
			}
			part[k][i] = val
		}
	}

	// DP for the first L parts
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, d+1)
	}
	dp[0][d] = 1
	for i := 1; i <= L; i++ {
		next := make([][]int, n+1)
		for x := range next {
			next[x] = make([]int, d+1)
		}
		bound := B[i-1]
		for s := 0; s <= bound; s++ {
			for prev := 1; prev <= d; prev++ {
				ways := dp[s][prev]
				if ways == 0 {
					continue
				}
				maxv := bound - s
				if maxv > prev {
					maxv = prev
				}
				if maxv > d {
					maxv = d
				}
				for v := 1; v <= maxv; v++ {
					ns := s + v
					next[ns][v] += ways
					if next[ns][v] >= MOD {
						next[ns][v] -= MOD
					}
				}
			}
		}
		dp = next
	}

	ans := 0
	for s := 0; s <= n; s++ {
		rem := n - s
		if rem < 0 {
			continue
		}
		for v := 1; v <= d; v++ {
			ways := dp[s][v]
			if ways == 0 {
				continue
			}
			ans = (ans + ways*part[v][rem]) % MOD
		}
	}
	fmt.Fprintln(out, ans%MOD)
}
