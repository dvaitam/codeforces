package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Query struct {
	l, r, t, idx int
}

type Update struct {
	pos, prev, val int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	orig := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &orig[i])
	}
	cur := make([]int, n+1)
	copy(cur, orig)

	values := make([]int, 0, n+q+5)
	values = append(values, orig[1:]...)

	queries := make([]Query, 0)
	updates := make([]Update, 0)

	for i := 0; i < q; i++ {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		if t == 1 {
			l := x
			r := y
			queries = append(queries, Query{l: l, r: r, t: len(updates), idx: len(queries)})
		} else {
			p := x
			val := y
			updates = append(updates, Update{pos: p, prev: cur[p], val: val})
			cur[p] = val
			values = append(values, val)
		}
	}

	// coordinate compression
	sort.Ints(values)
	uniq := make([]int, 0, len(values))
	for _, v := range values {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	mp := make(map[int]int, len(uniq))
	for i, v := range uniq {
		mp[v] = i
	}

	for i := 1; i <= n; i++ {
		orig[i] = mp[orig[i]]
	}
	for i := range updates {
		updates[i].prev = mp[updates[i].prev]
		updates[i].val = mp[updates[i].val]
	}

	arr := make([]int, n+1)
	copy(arr, orig)

	m := len(uniq)
	freq := make([]int, m)
	cntFreq := make([]int, n+2)
	cntFreq[0] = m
	mex := 1

	addVal := func(x int) {
		f := freq[x]
		cntFreq[f]--
		freq[x]++
		cntFreq[f+1]++
		for mex <= n && cntFreq[mex] > 0 {
			mex++
		}
	}
	removeVal := func(x int) {
		f := freq[x]
		cntFreq[f]--
		freq[x]--
		cntFreq[f-1]++
		for mex > 1 && cntFreq[mex-1] == 0 {
			mex--
		}
		for mex <= n && cntFreq[mex] > 0 {
			mex++
		}
	}

	// sort queries using Mo's algorithm with updates
	block := int(math.Pow(float64(n), 2.0/3.0)) + 1
	sort.Slice(queries, func(i, j int) bool {
		a, b := queries[i], queries[j]
		if a.l/block != b.l/block {
			return a.l < b.l
		}
		if a.r/block != b.r/block {
			return a.r < b.r
		}
		return a.t < b.t
	})

	ans := make([]int, len(queries))
	cl, cr, ct := 1, 0, 0
	for _, qu := range queries {
		for ct < qu.t {
			up := updates[ct]
			pos := up.pos
			if cl <= pos && pos <= cr {
				removeVal(arr[pos])
				addVal(up.val)
			}
			arr[pos] = up.val
			ct++
		}
		for ct > qu.t {
			up := updates[ct-1]
			pos := up.pos
			if cl <= pos && pos <= cr {
				removeVal(arr[pos])
				addVal(up.prev)
			}
			arr[pos] = up.prev
			ct--
		}
		for cl > qu.l {
			cl--
			addVal(arr[cl])
		}
		for cr < qu.r {
			cr++
			addVal(arr[cr])
		}
		for cl < qu.l {
			removeVal(arr[cl])
			cl++
		}
		for cr > qu.r {
			removeVal(arr[cr])
			cr--
		}
		for mex <= n && cntFreq[mex] > 0 {
			mex++
		}
		ans[qu.idx] = mex
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < len(ans); i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
