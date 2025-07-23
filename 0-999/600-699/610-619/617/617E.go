package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Query struct {
	l, r  int
	idx   int
	block int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ arr[i-1]
	}

	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		queries[i] = Query{l: l - 1, r: r, idx: i}
	}

	blockSize := int(math.Sqrt(float64(n))) + 1
	for i := range queries {
		queries[i].block = queries[i].l / blockSize
	}
	sort.Slice(queries, func(i, j int) bool {
		qi, qj := queries[i], queries[j]
		if qi.block != qj.block {
			return qi.block < qj.block
		}
		if qi.block%2 == 0 {
			return qi.r < qj.r
		}
		return qi.r > qj.r
	})

	// counts of prefix XOR values
	maxVal := 1 << 20
	cnt := make([]int, maxVal)

	ansArr := make([]int64, m)
	var curL, curR int = 0, -1
	var ans int64 = 0

	add := func(pos int) {
		val := pref[pos]
		if val^k < len(cnt) {
			ans += int64(cnt[val^k])
		}
		cnt[val]++
	}

	remove := func(pos int) {
		val := pref[pos]
		cnt[val]--
		if val^k < len(cnt) {
			ans -= int64(cnt[val^k])
		}
	}

	for _, q := range queries {
		L, R := q.l, q.r
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
		ansArr[q.idx] = ans
	}

	for i := 0; i < m; i++ {
		fmt.Fprintln(writer, ansArr[i])
	}
}
