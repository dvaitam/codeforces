package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	rd = bufio.NewReader(os.Stdin)
	wr = bufio.NewWriter(os.Stdout)
)

func readInt() int {
	n := 0
	sign := 1
	b, err := rd.ReadByte()
	for err == nil && (b < '0' || b > '9') {
		if b == '-' {
			sign = -1
		}
		b, err = rd.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		n = n*10 + int(b-'0')
		b, err = rd.ReadByte()
	}
	return n * sign
}

var (
	n, m    int
	head    []int
	to, nxt []int
	tp, num []int
	edgenum int
	s       []int
	ans     []int
	vis     []bool
)

func add(u, v, t, c int) {
	edgenum++
	to[edgenum] = v
	tp[edgenum] = t
	num[edgenum] = c
	nxt[edgenum] = head[u]
	head[u] = edgenum
}

func gg() {
	fmt.Fprintln(wr, "Impossible")
	wr.Flush()
	os.Exit(0)
}

func dfs(u, fa, E int) {
	vis[u] = true
	for i := head[u]; i != 0; i = nxt[i] {
		v := to[i]
		if !vis[v] {
			dfs(v, u, i)
		}
	}
	if E != 0 {
		eid := E >> 1
		ans[eid] = -s[u] * tp[E]
		s[fa] += s[u]
		s[u] = 0
	}
	if s[u] != 0 {
		gg()
	}
}

func main() {
	defer wr.Flush()
	n = readInt()
	s = make([]int, n+1)
	for i := 1; i <= n; i++ {
		s[i] = -readInt()
	}
	m = readInt()
	size := 2*m + 5
	head = make([]int, n+1)
	to = make([]int, size)
	nxt = make([]int, size)
	tp = make([]int, size)
	num = make([]int, size)
	ans = make([]int, m+1)
	vis = make([]bool, n+1)
	edgenum = 1
	for i := 1; i <= m; i++ {
		u := readInt()
		v := readInt()
		add(u, v, 1, i)
		add(v, u, -1, i)
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			dfs(i, 0, 0)
		}
	}
	for i := 1; i <= n; i++ {
		if s[i] != 0 {
			gg()
		}
	}
	fmt.Fprintln(wr, "Possible")
	for i := 1; i <= m; i++ {
		fmt.Fprintln(wr, ans[i])
	}
}
