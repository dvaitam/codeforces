package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct{ parent []int }

func NewDSU(n int) *DSU {
	p := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		p[i] = i
	}
	return &DSU{p}
}
func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}
func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x != y {
		d.parent[x] = y
	}
}

func main() {
	rd := bufio.NewReader(os.Stdin)
	wr := bufio.NewWriter(os.Stdout)
	defer wr.Flush()
	var t int
	fmt.Fscan(rd, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(rd, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(rd, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(rd, &b[i])
		}
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + a[i-1] - b[i-1]
		}
		if pref[n] != 0 {
			fmt.Fprintln(wr, "NO")
			for i := 0; i < m; i++ {
				var l, r int
				fmt.Fscan(rd, &l, &r)
			}
			continue
		}
		segs := make([][2]int, m)
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(rd, &l, &r)
			l--
			segs[i] = [2]int{l, r}
			adj[l] = append(adj[l], i)
			adj[r] = append(adj[r], i)
		}
		dsu := NewDSU(n)
		visited := make([]bool, n+1)
		queue := make([]int, 0)
		for i := 0; i <= n; i++ {
			if pref[i] == 0 {
				visited[i] = true
				queue = append(queue, i)
				dsu.Union(i, i+1)
			}
		}
		cnt := make([]int, m)
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for len(adj[v]) > 0 {
				id := adj[v][len(adj[v])-1]
				adj[v] = adj[v][:len(adj[v])-1]
				cnt[id]++
				if cnt[id] == 2 {
					l := segs[id][0]
					r := segs[id][1]
					if l > r {
						l, r = r, l
					}
					x := dsu.Find(l + 1)
					for x <= r {
						if !visited[x] {
							visited[x] = true
							queue = append(queue, x)
						}
						dsu.Union(x, x+1)
						x = dsu.Find(x)
					}
				}
			}
		}
		ok := true
		for i := 0; i <= n; i++ {
			if !visited[i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(wr, "YES")
		} else {
			fmt.Fprintln(wr, "NO")
		}
	}
}
