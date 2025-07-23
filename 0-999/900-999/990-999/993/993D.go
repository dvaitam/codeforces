package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int64, n)
	for i := range b {
		fmt.Fscan(reader, &b[i])
	}

	// sort indices by a descending
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		if a[idx[i]] == a[idx[j]] {
			return b[idx[i]] > b[idx[j]]
		}
		return a[idx[i]] > a[idx[j]]
	})

	sa := make([]int64, n)
	sb := make([]int64, n)
	for i, id := range idx {
		sa[i] = a[id]
		sb[i] = b[id]
	}

	// group indices by equal sa
	type group struct {
		start, end int // [start,end)
	}
	var groups []group
	for i := 0; i < n; {
		j := i
		for j < n && sa[j] == sa[i] {
			j++
		}
		groups = append(groups, group{start: i, end: j})
		i = j
	}

	maxA := sa[0]
	high := maxA * 1000
	low := int64(0)
	check := func(th int64) bool {
		inf := int64(1 << 60)
		dp := make([]int64, n+1)
		for i := range dp {
			dp[i] = inf
		}
		dp[0] = 0
		for _, g := range groups {
			m := g.end - g.start
			// sort group's indices by sb desc
			ids := make([]int, m)
			for i := 0; i < m; i++ {
				ids[i] = g.start + i
			}
			sort.Slice(ids, func(i, j int) bool {
				return sb[ids[i]] > sb[ids[j]]
			})
			prefix := make([]int64, m+1)
			for i := 0; i < m; i++ {
				prefix[i+1] = prefix[i] + sa[ids[i]]*1000 - th*sb[ids[i]]
			}
			newdp := make([]int64, n+1)
			for i := range newdp {
				newdp[i] = inf
			}
			for k := 0; k <= n; k++ {
				if dp[k] == inf {
					continue
				}
				minX := 0
				if m-k > 0 {
					minX = m - k
				}
				for x := minX; x <= m; x++ {
					newk := k + 2*x - m
					if newk < 0 || newk > n {
						continue
					}
					val := dp[k] + prefix[x]
					if val < newdp[newk] {
						newdp[newk] = val
					}
				}
			}
			dp = newdp
		}
		for k := int(n % 2); k <= n; k += 2 {
			if dp[k] <= 0 {
				return true
			}
		}
		return false
	}
	for low < high {
		mid := (low + high) / 2
		if check(mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	fmt.Fprintln(writer, low)
}
