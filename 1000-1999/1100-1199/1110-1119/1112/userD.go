package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func buildSA(s string) []int {
	n := len(s)
	sa := make([]int, n)
	rank := make([]int, n)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
		rank[i] = int(s[i])
	}
	k := 1
	for k < n {
		k2 := k
		sort.Slice(sa, func(i, j int) bool {
			if rank[sa[i]] != rank[sa[j]] {
				return rank[sa[i]] < rank[sa[j]]
			}
			ri, rj := -1, -1
			if sa[i]+k2 < n {
				ri = rank[sa[i]+k2]
			}
			if sa[j]+k2 < n {
				rj = rank[sa[j]+k2]
			}
			return ri < rj
		})
		tmp[sa[0]] = 0
		for i := 1; i < n; i++ {
			a1, b1 := rank[sa[i]], -1
			if sa[i]+k2 < n {
				b1 = rank[sa[i]+k2]
			}
			a2, b2 := rank[sa[i-1]], -1
			if sa[i-1]+k2 < n {
				b2 = rank[sa[i-1]+k2]
			}
			if a1 == a2 && b1 == b2 {
				tmp[sa[i]] = tmp[sa[i-1]]
			} else {
				tmp[sa[i]] = tmp[sa[i-1]] + 1
			}
		}
		copy(rank, tmp)
		if rank[sa[n-1]] == n-1 {
			break
		}
		k <<= 1
	}
	return sa
}

func buildLCP(s string, sa []int) []int {
	n := len(s)
	rank := make([]int, n)
	for i := 0; i < n; i++ {
		rank[sa[i]] = i
	}
	lcp := make([]int, n-1)
	k := 0
	for i := 0; i < n; i++ {
		r := rank[i]
		if r == n-1 {
			k = 0
			continue
		}
		j := sa[r+1]
		for i+k < n && j+k < n && s[i+k] == s[j+k] {
			k++
		}
		lcp[r] = k
		if k > 0 {
			k--
		}
	}
	return lcp
}

func calcMaxLen(s string, sa []int, lcp []int) []int {
	n := len(s)
	maxLen := make([]int, n)
	for idx := 0; idx < n; idx++ {
		pos := sa[idx]
		minL := int(1 << 30)
		for j := idx - 1; j >= 0; j-- {
			if lcp[j] < minL {
				minL = lcp[j]
			}
			if minL == 0 {
				break
			}
			if sa[j] < pos && minL > maxLen[pos] {
				maxLen[pos] = minL
			}
		}
		minL = int(1 << 30)
		for j := idx; j < n-1; j++ {
			if lcp[j] < minL {
				minL = lcp[j]
			}
			if minL == 0 {
				break
			}
			if sa[j+1] < pos && minL > maxLen[pos] {
				maxLen[pos] = minL
			}
		}
	}
	return maxLen
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b int
	fmt.Fscan(in, &n, &a, &b)
	var s string
	fmt.Fscan(in, &s)

	sa := buildSA(s)
	lcp := buildLCP(s, sa)
	maxLen := calcMaxLen(s, sa, lcp)

	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = 1 << 60
	}

	seen := [26]bool{}
	for i := 0; i < n; i++ {
		// cost for single character
		charCost := int64(a)
		if seen[s[i]-'a'] {
			charCost = min(charCost, int64(b))
		}
		if dp[i]+charCost < dp[i+1] {
			dp[i+1] = dp[i] + charCost
		}

		// cost for longer substrings
		L := maxLen[i]
		if L > n-i {
			L = n - i
		}
		if L >= 2 {
			val := dp[i] + int64(b)
			for l := 2; l <= L; l++ {
				if val < dp[i+l] {
					dp[i+l] = val
				}
			}
		}
		seen[s[i]-'a'] = true
	}

	fmt.Println(dp[n])
}