package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 30

type Basis struct {
	b [maxBit + 1]int
}

func (bs *Basis) Add(x int) {
	for i := maxBit; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			bs.b[i] = x
			return
		}
		if x^bs.b[i] < x {
			x ^= bs.b[i]
		} else {
			x ^= bs.b[i]
		}
	}
}

func (bs *Basis) Merge(o *Basis) {
	for i := maxBit; i >= 0; i-- {
		if o.b[i] != 0 {
			bs.Add(o.b[i])
		}
	}
}

func (bs *Basis) MaxXor() int {
	res := 0
	for i := maxBit; i >= 0; i-- {
		if bs.b[i] != 0 && (res^bs.b[i]) > res {
			res ^= bs.b[i]
		}
	}
	return res
}

func copyBasis(src *Basis) Basis {
	var d Basis
	d = *src
	return d
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		parent := make([]int, n+1)
		depth := make([]int, n+1)
		tin := make([]int, n+1)
		tout := make([]int, n+1)
		order := make([]int, 0, n)
		children := make([][]int, n+1)

		// iterative DFS to fill parent, depth, tin/tout and order
		type item struct{ v, p, idx int }
		stack := []item{{1, 0, 0}}
		time := 0
		for len(stack) > 0 {
			cur := &stack[len(stack)-1]
			v := cur.v
			if cur.idx == 0 {
				tin[v] = time
				time++
				parent[v] = cur.p
				if cur.p != 0 {
					depth[v] = depth[cur.p] + 1
					children[cur.p] = append(children[cur.p], v)
				}
			}
			if cur.idx < len(g[v]) {
				to := g[v][cur.idx]
				cur.idx++
				if to == cur.p {
					continue
				}
				stack = append(stack, item{to, v, 0})
			} else {
				tout[v] = time - 1
				order = append(order, v)
				stack = stack[:len(stack)-1]
			}
		}

		// compute down basis in postorder
		down := make([]Basis, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			down[v].Add(a[v])
			for _, to := range children[v] {
				down[v].Merge(&down[to])
			}
		}

		// compute up basis via preorder dfs
		up := make([]Basis, n+1)
		// process root's children
		type qitem struct{ v int }
		q := []qitem{{1}}
		for len(q) > 0 {
			cur := q[len(q)-1]
			q = q[:len(q)-1]
			v := cur.v
			m := len(children[v])
			pref := make([]Basis, m+1)
			suff := make([]Basis, m+1)
			for i := 0; i < m; i++ {
				pref[i+1] = copyBasis(&pref[i])
				pref[i+1].Merge(&down[children[v][i]])
			}
			for i := m - 1; i >= 0; i-- {
				suff[i] = copyBasis(&suff[i+1])
				suff[i].Merge(&down[children[v][i]])
			}
			for i, to := range children[v] {
				tmp := copyBasis(&up[v])
				tmp.Add(a[v])
				// merge siblings
				bs := copyBasis(&pref[i])
				bs.Merge(&suff[i+1])
				tmp.Merge(&bs)
				up[to] = tmp
				q = append(q, qitem{to})
			}
		}

		// binary lifting for ancestors
		LOG := 20
		upTable := make([][]int, LOG)
		for i := range upTable {
			upTable[i] = make([]int, n+1)
		}
		for i := 1; i <= n; i++ {
			upTable[0][i] = parent[i]
		}
		for k := 1; k < LOG; k++ {
			for i := 1; i <= n; i++ {
				upTable[k][i] = upTable[k-1][upTable[k-1][i]]
			}
		}

		isAncestor := func(u, v int) bool {
			return tin[u] <= tin[v] && tout[v] <= tout[u]
		}

		getKth := func(v, k int) int {
			for i := 0; i < LOG; i++ {
				if (k>>i)&1 == 1 {
					v = upTable[i][v]
				}
			}
			return v
		}

		totalBasis := down[1]

		var qCount int
		fmt.Fscan(reader, &qCount)
		for ; qCount > 0; qCount-- {
			var r, v int
			fmt.Fscan(reader, &r, &v)
			if r == v {
				fmt.Fprintln(writer, totalBasis.MaxXor())
				continue
			}
			if !isAncestor(v, r) {
				fmt.Fprintln(writer, down[v].MaxXor())
			} else {
				child := getKth(r, depth[r]-depth[v]-1)
				fmt.Fprintln(writer, up[child].MaxXor())
			}
		}
	}
}
