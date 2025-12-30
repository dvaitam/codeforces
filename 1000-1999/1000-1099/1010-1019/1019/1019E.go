package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Constants adjusted for the binarization process (node count can double)
const (
	MaxN = 400010      // Maximum number of nodes after binarization
	MaxM = MaxN << 1   // Maximum number of edges
)

// Global Variables
var (
	n, m, s, tt, cc int
	cnt             int
	res             int
	rt              int // Root edge index
	
	// Graph Arrays
	hd   [MaxN]int
	to   [MaxM]int
	nxt  [MaxM]int
	valA [MaxM]int
	valB [MaxM]int
	w    [MaxM]bool // Edge removed flag

	// Algorithm Arrays
	sz  [MaxN]int
	v   [MaxN]int // Temp array for DFS children
	
	// Edge reconstruction buffer
	g [MaxN]EdgeRec

	// Convex Hull Arrays
	P [MaxN]float64
	Q [MaxN]float64
	
	pArr [MaxN]Line
	qArr [MaxN]Line
	eArr []Line // Dynamic slice for the global hull lines
)

// Structs
type EdgeRec struct {
	x, y, a, b int
}

type Line struct {
	k, b int64
}

// Lines implements sort.Interface for []Line
type Lines []Line

func (a Lines) Len() int           { return len(a) }
func (a Lines) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Lines) Less(i, j int) bool {
	if a[i].k == a[j].k {
		return a[i].b > a[j].b // If slopes equal, larger intercept comes first (to be kept)
	}
	return a[i].k < a[j].k
}

// add inserts an undirected edge (u, v) with weights A, B
func add(x, y, A, B int) {
	tt++
	to[tt] = y
	nxt[tt] = hd[x]
	valA[tt] = A
	valB[tt] = B
	hd[x] = tt

	tt++
	to[tt] = x
	nxt[tt] = hd[y]
	valA[tt] = A
	valB[tt] = B
	hd[y] = tt
}

// chk determines if line b is redundant given lines a and c
// Intersection of (b, c) <= Intersection of (a, b)
func chk(a, b, c Line) bool {
	return float64(c.b-b.b)/float64(b.k-c.k) <= float64(b.b-a.b)/float64(a.k-b.k)
}

// dfs performs ternary binarization of the tree
// 
// Transforms a multi-child node into a chain of dummy nodes to ensure edge degree <= 3
func dfs(x, p int) {
	// Recurse first
	for i := hd[x]; i != 0; i = nxt[i] {
		if to[i] != p {
			dfs(to[i], x)
		}
	}

	// Collect children edges
	t := 0
	lst := x
	for i := hd[x]; i != 0; i = nxt[i] {
		if to[i] != p {
			t++
			v[t] = i
		}
	}

	// Rebuild edges into 'g' buffer
	if t > 0 {
		cc++
		g[cc] = EdgeRec{x, to[v[1]], valA[v[1]], valB[v[1]]}
		for i := 2; i < t; i++ {
			cc++
			n++ // Create dummy node
			g[cc] = EdgeRec{lst, n, 0, 0}
			lst = n
			cc++
			g[cc] = EdgeRec{lst, to[v[i]], valA[v[i]], valB[v[i]]}
		}
		if t > 1 {
			cc++
			g[cc] = EdgeRec{lst, to[v[t]], valA[v[t]], valB[v[t]]}
		}
	}
}

// gtrt finds the centroid edge to split the tree
func gtrt(x, p int) {
	sz[x] = 1
	for i := hd[x]; i != 0; i = nxt[i] {
		y := to[i]
		if y != p && !w[i] {
			gtrt(y, x)
			sz[x] += sz[y]
			tVal := sz[y]
			if s-sz[y] > tVal {
				tVal = s - sz[y]
			}
			if tVal < res {
				res = tVal
				rt = i
			}
		}
	}
}

// gtd collects all paths from node x in the current subtree
func gtd(x, p, t int, sArr []Line, tp *int, s1, s2 int64) {
	if t != 0 {
		*tp++
		sArr[*tp] = Line{s1, s2}
	}
	for i := hd[x]; i != 0; i = nxt[i] {
		y := to[i]
		if y != p && !w[i] {
			gtd(y, x, valA[i]+valB[i], sArr, tp, s1+int64(valA[i]), s2+int64(valB[i]))
		}
	}
}

// build constructs the upper convex hull from a set of lines
func build(sArr []Line, intersect []float64, tp *int, f bool) {
	// Sort lines by slope k
	slice := sArr[1 : *tp+1]
	sort.Sort(Lines(slice))

	t := 0
	// Graham scan-like algorithm to maintain upper envelope
	for i := 0; i < len(slice); i++ {
		ln := slice[i]
		if i == 0 || ln.k != sArr[t].k {
			for t > 1 && chk(sArr[t-1], sArr[t], ln) {
				t--
			}
			t++
			sArr[t] = ln
		}
	}
	*tp = t
	if f {
		for i := 1; i < *tp; i++ {
			intersect[i] = float64(sArr[i].b-sArr[i+1].b) / float64(sArr[i+1].k-sArr[i].k)
		}
	}
}

// sol is the main recursive solver (Edge Divide and Conquer)
func sol(x int) {
	if sz[x] == 1 {
		return
	}
	s = sz[x]
	res = 1 << 30
	gtrt(x, 0)

	// Mark edge as removed
	edgeIdx := rt
	w[edgeIdx] = true
	w[edgeIdx^1] = true
	
	nodeA := to[edgeIdx]
	nodeB := to[edgeIdx^1]

	// Update sizes for recursion
	if sz[nodeA] < sz[nodeB] {
		sz[nodeB] = s - sz[nodeA]
	} else {
		sz[nodeA] = s - sz[nodeB]
	}

	// Collect paths from both sides
	s1Cnt := 0
	s2Cnt := 0
	gtd(nodeA, nodeB, 1, pArr[:], &s1Cnt, int64(valA[edgeIdx]), int64(valB[edgeIdx]))
	gtd(nodeB, nodeA, 1, qArr[:], &s2Cnt, 0, 0)

	// Build hulls for merging
	build(pArr[:], P[:], &s1Cnt, true)
	build(qArr[:], Q[:], &s2Cnt, true)

	// Minkowski Sum-like merge of two convex hulls
	// Adds valid combined lines to the global list eArr
	j := 1
	for i := 1; i <= s1Cnt; i++ {
		for {
			cnt++
			if cnt >= len(eArr) {
				// Expand slice if needed
				newE := make([]Line, len(eArr)*2)
				copy(newE, eArr)
				eArr = newE
			}
			eArr[cnt] = Line{pArr[i].k + qArr[j].k, pArr[i].b + qArr[j].b}

			if j == s2Cnt {
				break
			}
			// Compare intersection points to advance pointers
			if i < s1Cnt && P[i] < Q[j] {
				break
			}
			j++
		}
	}

	sol(nodeA)
	sol(nodeB)
}

func main() {
	// Standard IO setup
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m)
	eArr = make([]Line, MaxN*20) // Initial allocation

	for i := 1; i < n; i++ {
		var u, v, x, y int
		fmt.Fscan(reader, &u, &v, &x, &y)
		add(u, v, x, y)
	}

	// 1. Transform tree
	dfs(1, 0)

	// 2. Rebuild graph with new edges
	tt = 1
	for i := 0; i < len(hd); i++ {
		hd[i] = 0
	}
	for i := 1; i <= cc; i++ {
		add(g[i].x, g[i].y, g[i].a, g[i].b)
	}

	// 3. Solve
	sz[1] = n
	sol(1)

	// 4. Build final global hull
	build(eArr, P[:], &cnt, false)

	// 5. Answer queries
	// Since queries x=0..m-1 are monotonic, we can just advance pointer j
	j := 1
	for i := 0; i < m; i++ {
		queryX := int64(i)
		for j < cnt && (eArr[j].k*queryX+eArr[j].b < eArr[j+1].k*queryX+eArr[j+1].b) {
			j++
		}
		fmt.Fprintf(writer, "%d ", eArr[j].k*queryX+eArr[j].b)
	}
}
