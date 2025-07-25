package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// query sends a query for the subset s and returns the number of edges inside s.
func query(s []int) int {
	// print query
	writer.WriteByte('?')
	writer.WriteByte(' ')
	writer.WriteString(fmt.Sprintf("%d", len(s)))
	for _, v := range s {
		writer.WriteByte(' ')
		writer.WriteString(fmt.Sprintf("%d", v))
	}
	writer.WriteByte('\n')
	writer.Flush()
	// read answer
	var ans int
	fmt.Fscan(reader, &ans)
	return ans
}

// findConflictPair finds two vertices u,v in group such that there is an edge between them.
func findConflictPair(group []int) (int, int) {
	var prefixE []int
	for i := range group {
		// edges in group[0..i]
		tot := query(group[:i+1])
		prefixE = append(prefixE, tot)
		prev := 0
		if i > 0 {
			prev = prefixE[i-1]
		}
		if tot > prev {
			u := group[i]
			// find partner v in group[0..i-1]
			l, r := 0, i-1
			for l < r {
				m := (l + r) / 2
				// check in subset group[0..m] with u
				sub := append(group[:m+1], u)
				tot2 := query(sub)
				if tot2 > prefixE[m] {
					r = m
				} else {
					l = m + 1
				}
			}
			v := group[l]
			return u, v
		}
	}
	return -1, -1
}

func main() {
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	// color: 0 or 1, default 0
	color := make([]int, n+1)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	// order of added vertices
	order := make([]int, 0, n)
	// prefixE[i] = edges in order[0..i]
	prefixE := make([]int, 0, n)
	// start with vertex 1
	order = append(order, 1)
	prefixE = append(prefixE, 0)
	parent[1] = 0
	depth[1] = 0
	color[1] = 0
	currE := 0
	// process rest
	for v := 2; v <= n; v++ {
		// edges in order ∪ {v}
		set := append(order, v)
		tot := query(set)
		delta := tot - currE
		if delta > 0 {
			// find one neighbor u in order
			l, r := 0, len(order)-1
			for l < r {
				m := (l + r) / 2
				sub := make([]int, m+2)
				copy(sub, order[:m+1])
				sub[m+1] = v
				tot2 := query(sub)
				if tot2-prefixE[m] > 0 {
					r = m
				} else {
					l = m + 1
				}
			}
			u := order[l]
			// assign v
			color[v] = color[u] ^ 1
			parent[v] = u
			depth[v] = depth[u] + 1
			// update order and prefixE
			order = append(order, v)
			prefixE = append(prefixE, tot)
			currE = tot
		} else {
			// no neighbors yet
			color[v] = 0
			parent[v] = 0
			depth[v] = 0
			order = append(order, v)
			prefixE = append(prefixE, currE)
		}
	}
	// check bipartiteness
	var grp [2][]int
	for i := 1; i <= n; i++ {
		grp[color[i]] = append(grp[color[i]], i)
	}
	// test each group for internal edge
	bad := -1
	for c := 0; c < 2; c++ {
		if len(grp[c]) > 0 {
			if query(grp[c]) > 0 {
				bad = c
				break
			}
		}
	}
	if bad < 0 {
		// bipartite: output partitions
		writer.WriteString("! Y\n")
		// first group
		writer.WriteString(fmt.Sprintf("%d", len(grp[0])))
		for _, v := range grp[0] {
			writer.WriteByte(' ')
			writer.WriteString(fmt.Sprintf("%d", v))
		}
		writer.WriteByte('\n')
		// second group
		writer.WriteString(fmt.Sprintf("%d", len(grp[1])))
		for _, v := range grp[1] {
			writer.WriteByte(' ')
			writer.WriteString(fmt.Sprintf("%d", v))
		}
		writer.WriteByte('\n')
		writer.Flush()
		return
	}
	// find an edge u-v inside bad group
	u, v := findConflictPair(grp[bad])
	// build cycle between u and v
	// ancestor paths
	au := map[int]bool{}
	x := u
	for x != 0 {
		au[x] = true
		x = parent[x]
	}
	// find LCA on path from v
	y := v
	for !au[y] {
		y = parent[y]
	}
	lca := y
	// build path u->lca
	pathU := []int{}
	x = u
	for x != lca {
		pathU = append(pathU, x)
		x = parent[x]
	}
	pathU = append(pathU, lca)
	// build path v->lca
	pathV := []int{}
	y = v
	for y != lca {
		pathV = append(pathV, y)
		y = parent[y]
	}
	// full cycle: pathU + reverse(pathV)
	cycle := make([]int, 0, len(pathU)+len(pathV))
	cycle = append(cycle, pathU...)
	for i := len(pathV) - 1; i >= 0; i-- {
		cycle = append(cycle, pathV[i])
	}
	// output cycle
	writer.WriteString("! N\n")
	writer.WriteString(fmt.Sprintf("%d", len(cycle)))
	for _, v := range cycle {
		writer.WriteByte(' ')
		writer.WriteString(fmt.Sprintf("%d", v))
	}
	writer.WriteByte('\n')
	writer.Flush()
}
