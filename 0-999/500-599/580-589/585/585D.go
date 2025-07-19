package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod  = 9999993
	base = 500000000
	inf  = 1000000000
)

type triple struct{ a, b, c int }
type entry struct{ a, b, c, sta int }

var (
	n, m       int
	p          []triple
	pow3       []int
	ans        int
	sta1, sta2 int
	buckets    map[int][]entry
)

func dfs1(k, lim, a, b, c, code int) {
	if k > lim {
		x := (a - b + base) % mod
		buckets[x] = append(buckets[x], entry{a, b, c, code})
		return
	}
	// choices: LM (a+=p[k].a, b+=p[k].b)
	dfs1(k+1, lim, a+p[k].a, b+p[k].b, c, code*3)
	// LW (a+=p[k].a, c+=p[k].c)
	dfs1(k+1, lim, a+p[k].a, b, c+p[k].c, code*3+1)
	// MW (b+=p[k].b, c+=p[k].c)
	dfs1(k+1, lim, a, b+p[k].b, c+p[k].c, code*3+2)
}

func dfs2(k, lim, a, b, c, code int) {
	if k > lim {
		// query
		x := (b - a + base) % mod
		tmpBest := -inf * 2
		var bestSta1 int
		if list, ok := buckets[x]; ok {
			for _, e := range list {
				if e.b-e.c == c-b {
					if e.a > tmpBest {
						tmpBest = e.a
						bestSta1 = e.sta
					}
				}
			}
		}
		if tmpBest > -inf {
			if tmpBest+a > ans {
				ans = tmpBest + a
				sta1 = code
				sta2 = bestSta1
			}
		}
		return
	}
	dfs2(k+1, lim, a+p[k].a, b+p[k].b, c, code*3)
	dfs2(k+1, lim, a+p[k].a, b, c+p[k].c, code*3+1)
	dfs2(k+1, lim, a, b+p[k].b, c+p[k].c, code*3+2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m = (n + 1) >> 1
	p = make([]triple, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i].a, &p[i].b, &p[i].c)
	}
	// precompute powers of 3
	pow3 = make([]int, n+2)
	pow3[0] = 1
	for i := 1; i <= n+1; i++ {
		pow3[i] = pow3[i-1] * 3
	}
	buckets = make(map[int][]entry)
	ans = -inf * 2
	// insert second half [m+1..n]
	dfs1(m+1, n, 0, 0, 0, 0)
	// query with first half [1..m]
	dfs2(1, m, 0, 0, 0, 0)
	if ans < -inf/2 {
		fmt.Fprintln(out, "Impossible")
		return
	}
	// output first half decisions
	for i := 1; i <= m; i++ {
		t := (sta1 / pow3[m-i]) % 3
		switch t {
		case 0:
			fmt.Fprintln(out, "LM")
		case 1:
			fmt.Fprintln(out, "LW")
		case 2:
			fmt.Fprintln(out, "MW")
		}
	}
	// output second half decisions
	for i := m + 1; i <= n; i++ {
		t := (sta2 / pow3[n-i]) % 3
		switch t {
		case 0:
			fmt.Fprintln(out, "LM")
		case 1:
			fmt.Fprintln(out, "LW")
		case 2:
			fmt.Fprintln(out, "MW")
		}
	}
}
