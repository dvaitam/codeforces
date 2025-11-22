package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type query struct {
	l, r, k, idx int
}

const (
	maxV      = 1 << 18 // 262144, covers XOR values up to 1.5e5
	bucketLen = 512
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		var q int
		fmt.Fscan(in, &q)
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &queries[i].l, &queries[i].r, &queries[i].k)
			queries[i].idx = i
		}

		blockSize := 1
		for blockSize*blockSize < n {
			blockSize++
		}

		sort.Slice(queries, func(i, j int) bool {
			bl1 := queries[i].l / blockSize
			bl2 := queries[j].l / blockSize
			if bl1 != bl2 {
				return bl1 < bl2
			}
			if bl1%2 == 0 {
				return queries[i].r < queries[j].r
			}
			return queries[i].r > queries[j].r
		})

		freq := make([]int, maxV)
		buckets := make([]int, (maxV+bucketLen-1)/bucketLen)

		updateFreq := func(val int, delta int) {
			freq[val] += delta
			buckets[val/bucketLen] += delta
		}

		kth := func(k int) int {
			k-- // zero-based
			bidx := 0
			for k >= buckets[bidx] {
				k -= buckets[bidx]
				bidx++
			}
			start := bidx * bucketLen
			for {
				if freq[start] > k {
					return start
				}
				k -= freq[start]
				start++
			}
		}

		inSet := make([]bool, n+1)
		val := make([]int, n+1)

		add := func(x int) {
			inSet[x] = true
			vx := 0
			for _, y := range adj[x] {
				if inSet[y] {
					vx ^= y
					old := val[y]
					newv := old ^ x
					updateFreq(old, -1)
					updateFreq(newv, 1)
					val[y] = newv
				}
			}
			val[x] = vx
			updateFreq(vx, 1)
		}

		remove := func(x int) {
			updateFreq(val[x], -1)
			inSet[x] = false
			for _, y := range adj[x] {
				if inSet[y] {
					old := val[y]
					newv := old ^ x
					updateFreq(old, -1)
					updateFreq(newv, 1)
					val[y] = newv
				}
			}
		}

		ans := make([]int, q)

		curL, curR := 1, 0 // empty interval
		for _, qu := range queries {
			L, R := qu.l, qu.r
			for curL > L {
				curL--
				add(curL)
			}
			for curR < R {
				curR++
				add(curR)
			}
			for curL < L {
				remove(curL)
				curL++
			}
			for curR > R {
				remove(curR)
				curR--
			}
			ans[qu.idx] = kth(qu.k)
		}

		for i := 0; i < q; i++ {
			if i > 0 {
				fmt.Fprint(out, "\n")
			}
			fmt.Fprint(out, ans[i])
		}
		if T > 0 {
			fmt.Fprint(out, "\n")
		}
	}
}
