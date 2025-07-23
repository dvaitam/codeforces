package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MAXN = 300005
	LOG  = 20
)

var depth [MAXN]int
var up [MAXN][LOG]int

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			u = up[u][i]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[u][i] != up[v][i] {
			u = up[u][i]
			v = up[v][i]
		}
	}
	return up[u][0]
}

func dist(u, v int) int {
	w := lca(u, v)
	return depth[u] + depth[v] - 2*depth[w]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}

	n := 1
	for i := 0; i < LOG; i++ {
		up[1][i] = 1
	}

	leaves := map[int]bool{1: true}
	A, B, D := 1, 1, 0
	res := make([]int, 0, m)
	count := 0

	for i := 0; i < m; i++ {
		var p int
		fmt.Fscan(reader, &p)
		n++
		depth[n] = depth[p] + 1
		up[n][0] = p
		for j := 1; j < LOG; j++ {
			up[n][j] = up[up[n][j-1]][j-1]
		}

		leaves[n] = true
		if leaves[p] {
			delete(leaves, p)
		}

		dA := dist(n, A)
		dB := dist(n, B)
		if dA > D && dA >= dB {
			B = n
			D = dA
		} else if dB > D {
			A = n
			D = dB
		}

		count = 0
		for v := range leaves {
			da := dist(v, A)
			db := dist(v, B)
			if da > db {
				if da == D {
					count++
				}
			} else {
				if db == D {
					count++
				}
			}
		}
		res = append(res, count)
	}

	writer := bufio.NewWriter(os.Stdout)
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
	writer.Flush()
}
