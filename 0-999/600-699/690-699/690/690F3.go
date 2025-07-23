package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// keyOf returns a unique key for slice c
func keyOf(c []int) string {
	if len(c) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(c[0]))
	for i := 1; i < len(c); i++ {
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(c[i]))
	}
	return sb.String()
}

// M stores canonical forms
var M map[string]int

// mapT maps slice c to unique int
func mapT(c []int) int {
	k := keyOf(c)
	if v, ok := M[k]; ok {
		return v
	}
	id := len(M)
	M[k] = id
	return id
}

// canonizeIndices remaps vertex labels to 0..n-1 by first appearance
func canonizeIndices(e [][2]int) [][2]int {
	m2 := make(map[int]int)
	idx := 0
	for i := range e {
		u := e[i][0]
		v := e[i][1]
		if _, ok := m2[u]; !ok {
			m2[u] = idx
			idx++
		}
		if _, ok := m2[v]; !ok {
			m2[v] = idx
			idx++
		}
		e[i][0] = m2[u]
		e[i][1] = m2[v]
	}
	return e
}

// encodeSubtree returns canonical label of subtree rooted at x
func encodeSubtree(x, dad int, E [][]int) int {
	var ch []int
	for _, y := range E[x] {
		if y != dad {
			ch = append(ch, encodeSubtree(y, x, E))
		}
	}
	sort.Ints(ch)
	return mapT(ch)
}

func dfsC(x, dad int, E [][]int, sz, bal []int) {
	sz[x] = 1
	bal[x] = 0
	for _, y := range E[x] {
		if y != dad {
			dfsC(y, x, E, sz, bal)
			if sz[y] > bal[x] {
				bal[x] = sz[y]
			}
			sz[x] += sz[y]
		}
	}
}

// canonizeTree returns canonical label of tree e
func canonizeTree(e [][2]int) int {
	e = canonizeIndices(e)
	n := 0
	for _, p := range e {
		if p[0]+1 > n {
			n = p[0] + 1
		}
		if p[1]+1 > n {
			n = p[1] + 1
		}
	}
	E := make([][]int, n)
	for _, p := range e {
		E[p[0]] = append(E[p[0]], p[1])
		E[p[1]] = append(E[p[1]], p[0])
	}
	sz := make([]int, n)
	bal := make([]int, n)
	dfsC(0, -1, E, sz, bal)
	ans := -1
	for x := 0; x < n; x++ {
		rem := n - sz[x]
		if rem > bal[x] {
			bal[x] = rem
		}
		if 2*bal[x] <= n {
			h := encodeSubtree(x, -1, E)
			if ans == -1 || h < ans {
				ans = h
			}
		}
	}
	return ans
}

func dfsMark(x, t int, E [][]int, bio []int) {
	bio[x] = t
	for _, y := range E[x] {
		if bio[y] == 0 {
			dfsMark(y, t, E, bio)
		}
	}
}

// canonizeForest returns canonical label of forest e
func canonizeForest(e [][2]int) int {
	e = canonizeIndices(e)
	n := 0
	for _, p := range e {
		if p[0]+1 > n {
			n = p[0] + 1
		}
		if p[1]+1 > n {
			n = p[1] + 1
		}
	}
	E := make([][]int, n)
	for _, p := range e {
		E[p[0]] = append(E[p[0]], p[1])
		E[p[1]] = append(E[p[1]], p[0])
	}
	bio := make([]int, n)
	var all []int
	for i := 0; i < n; i++ {
		if bio[i] == 0 {
			dfsMark(i, i+1, E, bio)
			var ne [][2]int
			for _, p := range e {
				if bio[p[0]] == i+1 && bio[p[1]] == i+1 {
					ne = append(ne, p)
				}
			}
			all = append(all, canonizeTree(ne))
		}
	}
	sort.Ints(all)
	return mapT(all)
}

// attempt tries to reconstruct the tree assuming base is the drawing
// where a leaf was removed. target is the canonical hash of the other drawing.
func attempt(base [][2]int, target int, n int) [][2]int {
	deg := (n - 1) - len(base)
	if deg != 1 {
		return nil
	}
	m := 0
	for _, p := range base {
		if p[0]+1 > m {
			m = p[0] + 1
		}
		if p[1]+1 > m {
			m = p[1] + 1
		}
	}
	newNode := m
	original := make([][2]int, len(base))
	copy(original, base)
	for i := 0; i < newNode; i++ {
		g := append(original, [2]int{i, newNode})
		for j := 0; j <= newNode; j++ {
			var ng [][2]int
			for _, p := range g {
				if p[0] != j && p[1] != j {
					ng = append(ng, p)
				}
			}
			if canonizeForest(ng) == target {
				ans := make([][2]int, len(g))
				copy(ans, g)
				return ans
			}
		}
	}
	return nil
}

func solve() [][2]int {
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return nil
	}
	v := make([][][2]int, k)
	for i := 0; i < k; i++ {
		var m int
		fmt.Fscan(reader, &m)
		v[i] = make([][2]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &v[i][j][0], &v[i][j][1])
			v[i][j][0]--
			v[i][j][1]--
		}
		v[i] = canonizeIndices(v[i])
	}
	M = make(map[string]int)
	h0 := canonizeForest(v[0])
	h1 := canonizeForest(v[1])
	if ans := attempt(v[0], h1, n); ans != nil {
		return ans
	}
	if ans := attempt(v[1], h0, n); ans != nil {
		return ans
	}
	return nil
}

func main() {
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		ans := solve()
		if ans == nil {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
			for _, p := range ans {
				fmt.Fprintf(writer, "%d %d\n", p[0]+1, p[1]+1)
			}
		}
	}
}
