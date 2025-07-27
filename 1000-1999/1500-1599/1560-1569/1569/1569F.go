package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct{ a, b int }

type DSU struct{ parent []int }

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func encodeDSU(d *DSU) string {
	n := len(d.parent)
	mp := make(map[int]byte)
	next := byte(0)
	arr := make([]byte, n)
	for i := 0; i < n; i++ {
		r := d.find(i)
		v, ok := mp[r]
		if !ok {
			v = next
			mp[r] = v
			next++
		}
		arr[i] = v
	}
	return string(arr)
}

var (
	n, m, k int
	adj     [][]bool
	dsuSet  map[string]struct{}
)

func feasible(pairs []Pair) bool {
	m := len(pairs)
	full := (1 << m) - 1
	memo := make([][][]int8, 1<<m)
	for i := range memo {
		memo[i] = make([][]int8, n)
		for j := 0; j < n; j++ {
			memo[i][j] = make([]int8, n)
			for l := 0; l < n; l++ {
				memo[i][j][l] = -1
			}
		}
	}
	var dfs func(mask int, left, right int) bool
	dfs = func(mask int, left, right int) bool {
		if mask == full {
			return true
		}
		if memo[mask][left][right] != -1 {
			return memo[mask][left][right] == 1
		}
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				continue
			}
			a := pairs[i].a
			b := pairs[i].b
			if adj[left][a] && adj[right][b] {
				if dfs(mask|1<<i, a, b) {
					memo[mask][left][right] = 1
					return true
				}
			}
			if adj[left][b] && adj[right][a] {
				if dfs(mask|1<<i, b, a) {
					memo[mask][left][right] = 1
					return true
				}
			}
		}
		memo[mask][left][right] = 0
		return false
	}

	for i := 0; i < m; i++ {
		a, b := pairs[i].a, pairs[i].b
		if !adj[a][b] {
			continue
		}
		if dfs(1<<i, a, b) {
			return true
		}
		if dfs(1<<i, b, a) {
			return true
		}
	}
	return false
}

func addCoarsenings(pairs []Pair) {
	m := len(pairs)
	group := make([]int, m)
	var dfs func(idx, used int)
	dfs = func(idx, used int) {
		if idx == m {
			d := NewDSU(n)
			for i := 0; i < m; i++ {
				d.union(pairs[i].a, pairs[i].b)
			}
			for g := 0; g < used; g++ {
				first := -1
				for i := 0; i < m; i++ {
					if group[i] == g {
						if first == -1 {
							first = pairs[i].a
						} else {
							d.union(first, pairs[i].a)
						}
					}
				}
			}
			key := encodeDSU(d)
			dsuSet[key] = struct{}{}
			return
		}
		for g := 0; g < used; g++ {
			group[idx] = g
			dfs(idx+1, used)
		}
		group[idx] = used
		dfs(idx+1, used+1)
	}
	dfs(0, 0)
}

func genPairs(start int, used []bool, cur []Pair) {
	for start < n && used[start] {
		start++
	}
	if start == n {
		if feasible(cur) {
			addCoarsenings(cur)
		}
		return
	}
	used[start] = true
	for j := start + 1; j < n; j++ {
		if !used[j] {
			used[j] = true
			cur = append(cur, Pair{start, j})
			genPairs(start+1, used, cur)
			cur = cur[:len(cur)-1]
			used[j] = false
		}
	}
	used[start] = false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	adj = make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u][v] = true
		adj[v][u] = true
	}
	dsuSet = make(map[string]struct{})
	used := make([]bool, n)
	genPairs(0, used, nil)
	powk := make([]int64, n+1)
	powk[0] = 1
	for i := 1; i <= n; i++ {
		powk[i] = powk[i-1] * int64(k)
	}
	var ans int64 = 0
	for key := range dsuSet {
		maxc := 0
		for i := 0; i < len(key); i++ {
			if int(key[i])+1 > maxc {
				maxc = int(key[i]) + 1
			}
		}
		ans += powk[maxc]
	}
	fmt.Println(ans)
}
