package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxn = 1000009

var (
	a       = make([]int, maxn)
	res     = make([]int, maxn)
	cnt     = make([]int, maxn)
	f       = make([]int, maxn)
	p       = make([]int, maxn)
	q       = make([]int, maxn)
	g       = make([][]int, maxn)
	queries = make([][]qry, maxn)
)

type qry struct {
	i, l, k int
}

func swapPos(pos, newpos int) {
	x := p[pos]
	y := p[newpos]
	p[pos], p[newpos] = y, x
	q[y], q[x] = pos, newpos
}

func insert(x int) {
	cnt[x]++
	f[cnt[x]]++
	newpos := f[cnt[x]]
	pos := q[x]
	swapPos(pos, newpos)
}

func erase(x int) {
	pos := q[x]
	newpos := f[cnt[x]]
	f[cnt[x]]--
	cnt[x]--
	swapPos(pos, newpos)
}

func query(l, k int) int {
	if f[l] < k {
		return -1
	}
	return p[f[l]-k+1]
}

func solution(reader *bufio.Reader, writer *bufio.Writer) {
	var n, que int
	fmt.Fscan(reader, &n, &que)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		g[i] = g[i][:0]
		queries[i] = queries[i][:0]
	}
	for v := 2; v <= n; v++ {
		var u int
		fmt.Fscan(reader, &u)
		g[u] = append(g[u], v)
	}
	for i := 0; i <= n; i++ {
		p[i] = i
		q[i] = i
	}
	for i := 0; i < que; i++ {
		var u, l, k int
		fmt.Fscan(reader, &u, &l, &k)
		queries[u] = append(queries[u], qry{i, l, k})
	}
	var dfs func(u int)
	dfs = func(u int) {
		insert(a[u])
		for _, qq := range queries[u] {
			res[qq.i] = query(qq.l, qq.k)
		}
		for _, v := range g[u] {
			dfs(v)
		}
		erase(a[u])
	}
	dfs(1)
	for i := 0; i < que; i++ {
		fmt.Fprint(writer, res[i], " ")
	}
	fmt.Fprintln(writer)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var tc int
	fmt.Fscan(reader, &tc)
	for tc > 0 {
		solution(reader, writer)
		tc--
	}
}
