package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 250001

var (
	n   int
	adj [MAXN][]int
	par [MAXN][18]int
	st  [MAXN]int
	en  [MAXN]int
	ptr int
	dsu [MAXN]int
	aw  []int
	ax  []int
	ay  []int
	az  []int
	t2  [MAXN][]int
	die [MAXN]bool
	deg [MAXN]int
)

func dfs(u, p int) {
	ptr++
	st[u] = ptr
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		par[v][0] = u
		for i := 1; i < 18; i++ {
			par[v][i] = par[par[v][i-1]][i-1]
		}
		dfs(v, u)
	}
	en[u] = ptr
}

func find(u int) int {
	if dsu[u] != u {
		dsu[u] = find(dsu[u])
	}
	return dsu[u]
}

func join(u, v int) {
	u = find(u)
	v = find(v)
	dsu[u] = v
}

func get(u, r int) int {
	for i := 17; i >= 0; i-- {
		pu := par[u][i]
		if pu == 0 {
			continue
		}
		du := find(pu)
		if st[r] < st[du] && en[du] <= en[r] {
			u = pu
		}
	}
	return u
}

func del(u, v int) {
	aw = append(aw, u)
	ax = append(ax, v)
	ru := find(u)
	rv := find(v)
	if st[ru] <= st[rv] && en[ru] >= en[rv] {
		w := get(rv, ru)
		ay = append(ay, w)
		az = append(az, par[w][0])
		join(w, par[w][0])
	} else {
		ay = append(ay, ru)
		az = append(az, par[ru][0])
		join(ru, par[ru][0])
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dfs(1, 0)
	for i := 1; i <= n; i++ {
		dsu[i] = i
	}
	aw = make([]int, 0, n-1)
	ax = make([]int, 0, n-1)
	ay = make([]int, 0, n-1)
	az = make([]int, 0, n-1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		deg[u]++
		deg[v]++
		t2[u] = append(t2[u], v)
		t2[v] = append(t2[v], u)
	}
	// initialize stack of leaves
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			stack = append(stack, i)
		}
	}
	for i := 1; i < n; i++ {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		die[u] = true
		var v int
		for _, x := range t2[u] {
			if !die[x] {
				v = x
				break
			}
		}
		del(u, v)
		deg[v]--
		if deg[v] == 1 {
			stack = append(stack, v)
		}
	}
	fmt.Fprintln(writer, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fprintf(writer, "%d %d %d %d\n", ay[i], az[i], aw[i], ax[i])
	}
}
