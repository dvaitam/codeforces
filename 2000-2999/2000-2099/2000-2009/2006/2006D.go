package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fenwick struct {
	n int
	f []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, f: make([]int, n+1)}
}

func (b *fenwick) add(idx, delta int) {
	for idx <= b.n {
		b.f[idx] += delta
		idx += idx & -idx
	}
}

func (b *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b.f[idx]
		idx -= idx & -idx
	}
	return res
}

type rmq struct {
	log []int
	st  [][]int
}

func buildRMQ(a []int) *rmq {
	n := len(a)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i>>1] + 1
	}
	k := log[n] + 1
	st := make([][]int, k)
	st[0] = make([]int, n)
	copy(st[0], a)
	for j := 1; j < k; j++ {
		st[j] = make([]int, n-(1<<j)+1)
		for i := 0; i+(1<<j) <= n; i++ {
			left := st[j-1][i]
			right := st[j-1][i+(1<<(j-1))]
			if right > left {
				left = right
			}
			st[j][i] = left
		}
	}
	return &rmq{log: log, st: st}
}

func (r *rmq) query(l, rr int) int {
	// l and rr are 0-indexed inclusive
	length := rr - l + 1
	k := r.log[length]
	a := r.st[k][l]
	b := r.st[k][rr-(1<<k)+1]
	if b > a {
		return b
	}
	return a
}

type qInfo struct {
	l, r   int
	limit  int
	idx    int
	length int
	safe   int
}

type valPos struct {
	val int
	pos int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q, k int
		fmt.Fscan(in, &n, &q, &k)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		rmqMax := buildRMQ(b)

		queries := make([]qInfo, q)
		for i := 0; i < q; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			r--
			length := r - l + 1
			vmax := rmqMax.query(l, r)
			limit := k / vmax
			queries[i] = qInfo{l: l, r: r, limit: limit, idx: i, length: length}
		}

		// offline counting of elements <= limit in ranges
		vals := make([]valPos, n)
		for i, v := range b {
			vals[i] = valPos{val: v, pos: i + 1} // fenwick uses 1-indexed positions
		}
		sort.Slice(vals, func(i, j int) bool {
			return vals[i].val < vals[j].val
		})

		sortedQueries := make([]qInfo, q)
		copy(sortedQueries, queries)
		sort.Slice(sortedQueries, func(i, j int) bool {
			return sortedQueries[i].limit < sortedQueries[j].limit
		})

		fw := newFenwick(n)
		ptr := 0
		for _, qu := range sortedQueries {
			for ptr < n && vals[ptr].val <= qu.limit {
				fw.add(vals[ptr].pos, 1)
				ptr++
			}
			safe := fw.sum(qu.r+1) - fw.sum(qu.l)
			queries[qu.idx].safe = safe
		}

		ans := make([]int, q)
		for _, qu := range queries {
			safe := qu.safe
			heavy := qu.length - safe
			if heavy <= 0 {
				ans[qu.idx] = 0
				continue
			}
			needSafe := heavy - 1
			if safe >= needSafe {
				ans[qu.idx] = 0
				continue
			}
			missing := needSafe - safe
			ans[qu.idx] = (missing + 1) / 2
		}

		for i := 0; i < q; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
