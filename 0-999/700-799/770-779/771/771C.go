package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, k       int
	adj        [][]int
	removed    []bool
	size       []int
	totalFloor int64
	totalDiv0  int64
)

func dfsSize(u, p int) int {
	size[u] = 1
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			size[u] += dfsSize(v, u)
		}
	}
	return size[u]
}

func findCentroid(u, p, tot int) int {
	for _, v := range adj[u] {
		if v != p && !removed[v] && size[v] > tot/2 {
			return findCentroid(v, u, tot)
		}
	}
	return u
}

func collect(u, p, dist int, cnt, sum []int64) {
	r := dist % k
	cnt[r]++
	sum[r] += int64(dist / k)
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			collect(v, u, dist+1, cnt, sum)
		}
	}
}

func decompose(u int) {
	tot := dfsSize(u, -1)
	c := findCentroid(u, -1, tot)
	removed[c] = true

	globalCnt := make([]int64, k)
	globalSum := make([]int64, k)
	globalCnt[0] = 1 // centroid itself

	for _, v := range adj[c] {
		if removed[v] {
			continue
		}
		cnt := make([]int64, k)
		sum := make([]int64, k)
		collect(v, c, 1, cnt, sum)
		for r := 0; r < k; r++ {
			if cnt[r] == 0 {
				continue
			}
			for b := 0; b < k; b++ {
				if globalCnt[b] == 0 {
					continue
				}
				carry := int64(0)
				if r+b >= k {
					carry = 1
				}
				pairCnt := cnt[r] * globalCnt[b]
				totalFloor += cnt[r]*globalSum[b] + globalCnt[b]*sum[r] + carry*pairCnt
			}
			totalDiv0 += cnt[r] * globalCnt[(k-r)%k]
		}
		for i := 0; i < k; i++ {
			globalCnt[i] += cnt[i]
			globalSum[i] += sum[i]
		}
	}

	for _, v := range adj[c] {
		if !removed[v] {
			decompose(v)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &k)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	removed = make([]bool, n+1)
	size = make([]int, n+1)
	decompose(1)
	totalPairs := int64(n) * int64(n-1) / 2
	ans := totalPairs + totalFloor - totalDiv0
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, ans)
	writer.Flush()
}
