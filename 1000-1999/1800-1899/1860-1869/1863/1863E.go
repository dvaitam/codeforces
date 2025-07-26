package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		var k int64
		fmt.Fscan(reader, &n, &m, &k)
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		adj := make([][]int, n)
		indeg := make([]int, n)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			a--
			b--
			adj[a] = append(adj[a], b)
			indeg[b]++
		}
		q := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				q = append(q, i)
			}
		}
		for idx := 0; idx < len(q); idx++ {
			u := q[idx]
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					q = append(q, v)
				}
			}
		}
		dp := make([]int64, n)
		for i := n - 1; i >= 0; i-- {
			u := q[i]
			for _, v := range adj[u] {
				w := (h[v] - h[u]) % k
				if w < 0 {
					w += k
				}
				if dp[u] < dp[v]+w {
					dp[u] = dp[v] + w
				}
			}
		}
		indeg = make([]int, n)
		for u := 0; u < n; u++ {
			for _, v := range adj[u] {
				indeg[v]++
			}
		}
		sources := []int{}
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				sources = append(sources, i)
			}
		}
		s := len(sources)
		if s == 0 {
			fmt.Println(0)
			continue
		}
		hs := make([]int64, s)
		ls := make([]int64, s)
		for i, idx := range sources {
			hs[i] = h[idx]
			ls[i] = dp[idx]
		}
		order := make([]int, s)
		for i := 0; i < s; i++ {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool { return hs[order[i]] < hs[order[j]] })
		hsorted := make([]int64, s)
		lsorted := make([]int64, s)
		for i, idx := range order {
			hsorted[i] = hs[idx]
			lsorted[i] = ls[idx]
		}
		prefMaxHL := make([]int64, s)
		prefMinH := make([]int64, s)
		var mx int64 = -1 << 60
		var mn int64 = 1 << 60
		for i := 0; i < s; i++ {
			val := hsorted[i] + lsorted[i]
			if val > mx {
				mx = val
			}
			prefMaxHL[i] = mx
			if hsorted[i] < mn {
				mn = hsorted[i]
			}
			prefMinH[i] = mn
		}
		suffMaxHL := make([]int64, s)
		suffMinH := make([]int64, s)
		mx = -1 << 60
		mn = 1 << 60
		for i := s - 1; i >= 0; i-- {
			val := hsorted[i] + lsorted[i]
			if val > mx {
				mx = val
			}
			suffMaxHL[i] = mx
			if hsorted[i] < mn {
				mn = hsorted[i]
			}
			suffMinH[i] = mn
		}
		unique := make(map[int64]bool)
		res := int64(1 << 62)
		for _, c := range hsorted {
			if unique[c] {
				continue
			}
			unique[c] = true
			idx := sort.Search(len(hsorted), func(i int) bool { return hsorted[i] >= c })
			max1 := int64(-1 << 60)
			min1 := int64(1 << 60)
			if idx < s {
				max1 = suffMaxHL[idx] - c
				min1 = suffMinH[idx] - c
			}
			max2 := int64(-1 << 60)
			min2 := int64(1 << 60)
			if idx > 0 {
				max2 = prefMaxHL[idx-1] - c + k
				min2 = prefMinH[idx-1] - c + k
			}
			maxv := max1
			if max2 > maxv {
				maxv = max2
			}
			minv := min1
			if min2 < minv {
				minv = min2
			}
			diff := maxv - minv
			if diff < res {
				res = diff
			}
		}
		fmt.Println(res)
	}
}
