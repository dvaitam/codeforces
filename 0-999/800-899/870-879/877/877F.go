package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Query struct {
	L   int
	R   int
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &t[i])
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		val := a[i-1]
		if t[i-1] == 2 {
			val = -val
		}
		prefix[i] = prefix[i-1] + val
	}

	var q int
	fmt.Fscan(reader, &q)
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		queries[i] = Query{L: l - 1, R: r, idx: i}
	}

	blockSize := int(math.Sqrt(float64(n+1))) + 1
	sort.Slice(queries, func(i, j int) bool {
		bi := queries[i].L / blockSize
		bj := queries[j].L / blockSize
		if bi != bj {
			return queries[i].L < queries[j].L
		}
		if bi%2 == 0 {
			return queries[i].R < queries[j].R
		}
		return queries[i].R > queries[j].R
	})

	freq := make(map[int64]int64)
	curL, curR := 0, -1
	var curAns int64
	res := make([]int64, q)

	addRight := func(pos int) {
		val := prefix[pos]
		curAns += freq[val-k]
		freq[val]++
	}
	removeRight := func(pos int) {
		val := prefix[pos]
		freq[val]--
		curAns -= freq[val-k]
		if freq[val] == 0 {
			delete(freq, val)
		}
	}
	addLeft := func(pos int) {
		val := prefix[pos]
		curAns += freq[val+k]
		freq[val]++
	}
	removeLeft := func(pos int) {
		val := prefix[pos]
		freq[val]--
		curAns -= freq[val+k]
		if freq[val] == 0 {
			delete(freq, val)
		}
	}

	for _, qu := range queries {
		L, R := qu.L, qu.R
		for curR < R {
			curR++
			addRight(curR)
		}
		for curR > R {
			removeRight(curR)
			curR--
		}
		for curL < L {
			removeLeft(curL)
			curL++
		}
		for curL > L {
			curL--
			addLeft(curL)
		}
		res[qu.idx] = curAns
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, res[i])
	}
}
