package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	N = 300000
	M = 300000
)

var (
	X      uint64
	ii     [M]int
	jj     [M]int
	eh     [][]int
	pp     [N]int
	ta     [N]int
	tb     [N]int
	prev_  [N]int
	next_  [N]int
	xx     [N]int
	inArr  [N]bool
	tree   [M]bool
	ds     [M]int
	uCount int
)

func srand_() {
	X = uint64(time.Now().UnixNano()) | 1
}

func rand_() int {
	X *= 3
	return int(X >> 1)
}

func appendEdge(i, h int) {
	eh[i] = append(eh[i], h)
}

func findSet(i int) int {
	if ds[i] < 0 {
		return i
	}
	ds[i] = findSet(ds[i])
	return ds[i]
}

func join(i, j int) {
	i = findSet(i)
	j = findSet(j)
	if i == j {
		return
	}
	if ds[i] > ds[j] {
		ds[i] = j
	} else {
		if ds[i] == ds[j] {
			ds[j] = ds[j] - 1
		}
		ds[j] = i
	}
}

func dfs(p, f, i int) {
	uCount++
	pp[i] = p
	ta[i] = uCount
	tb[i] = uCount
	for _, h := range eh[i] {
		j := i ^ ii[h] ^ jj[h]
		if j == p {
			continue
		}
		if ta[j] == 0 {
			tree[h] = true
			dfs(i, h, j)
			if tb[j] < tb[i] {
				tb[i] = tb[j]
			}
			if tb[j] < ta[i] {
				join(h, f)
			}
		} else if ta[j] < ta[i] {
			if tb[i] > ta[j] {
				tb[i] = ta[j]
			}
			join(h, f)
		}
	}
}

func compareComp(h1, h2 int) int {
	r1 := findSet(h1)
	r2 := findSet(h2)
	if r1 != r2 {
		return r1 - r2
	}
	return ta[ii[h1]] - ta[ii[h2]]
}

func compareLR(h1, h2 int) int {
	if xx[ii[h1]] != xx[ii[h2]] {
		return xx[ii[h1]] - xx[ii[h2]]
	}
	return xx[jj[h2]] - xx[jj[h1]]
}

// custom quicksort
func sortEdges(hh []int, l, r int, cmp func(int, int) int) {
	for l < r {
		i, j, k := l, l, r
		pivot := hh[l+rand_()%(r-l)]
		for j < k {
			c := cmp(hh[j], pivot)
			if c == 0 {
				j++
			} else if c < 0 {
				hh[i], hh[j] = hh[j], hh[i]
				i++
				j++
			} else {
				k--
				hh[j], hh[k] = hh[k], hh[j]
			}
		}
		sortEdges(hh, l, i, cmp)
		l = k
	}
}

func solve(hh []int, m int) bool {
	qu := make([]int, 0, N)
	h := hh[0]
	i := ii[h]
	j := jj[h]
	if m == 1 {
		eh[i] = append(eh[i], h)
		eh[j] = append(eh[j], h)
		return true
	}
	qu = append(qu, i)
	inArr[i] = true
	prev_[j] = i
	next_[i] = j
	for j != i {
		inArr[j] = true
		qu = append(qu, j)
		next_[pp[j]] = j
		prev_[j] = pp[j]
		j = pp[j]
	}
	for _, h = range hh[1:m] {
		i = ii[h]
		j = jj[h]
		if tree[h] || inArr[j] {
			continue
		}
		for !inArr[j] {
			j = pp[j]
		}
		if j == prev_[i] {
			k := i
			for !inArr[j] {
				inArr[j] = true
				qu = append(qu, j)
				next_[prev_[k]] = j
				prev_[j] = prev_[k]
				next_[j] = k
				k = j
				j = pp[j]
			}
		} else if j == next_[i] {
			k := i
			for !inArr[j] {
				inArr[j] = true
				qu = append(qu, j)
				next_[j] = next_[k]
				prev_[next_[k]] = j
				next_[k] = j
				prev_[j] = k
				k = j
				j = pp[j]
			}
		} else {
			for _, v := range qu {
				inArr[v] = false
			}
			return false
		}
	}
	// order cycle
	start := ii[hh[0]]
	cnt := 0
	j = start
	cycle := []int{}
	for {
		cycle = append(cycle, j)
		j = next_[j]
		if j == start {
			break
		}
	}
	for idx, v := range cycle {
		inArr[v] = false
		xx[v] = idx
	}
	for idx, h = range hh[:m] {
		if xx[ii[h]] > xx[jj[h]] {
			ii[h], jj[h] = jj[h], ii[h]
		}
	}
	sortEdges(hh, 0, m, compareLR)
	// check intervals
	stack := []int{}
	for _, h = range hh[:m] {
		for len(stack) > 0 && xx[jj[stack[len(stack)-1]]] <= xx[ii[h]] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 && xx[jj[h]] > xx[jj[stack[len(stack)-1]]] {
			return false
		}
		stack = append(stack, h)
	}
	// add edges
	for _, h = range hh[:m] {
		eh[ii[h]] = append(eh[ii[h]], h)
	}
	for _, h = range hh[:m] {
		eh[jj[h]] = append(eh[jj[h]], h)
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	srand_()
	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		var n, m int
		fmt.Fscan(reader, &n, &m)
		// init
		eh = make([][]int, n)
		for i := 0; i < n; i++ {
			eh[i] = eh[i][:0]
			prev_[i] = -1
			next_[i] = -1
		}
		for i := 0; i < m; i++ {
			ds[i] = -1
			tree[i] = false
		}
		for h := 0; h < m; h++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			ii[h] = u
			jj[h] = v
			appendEdge(u, h)
			appendEdge(v, h)
		}
		// biconnected check
		for i := 0; i < n; i++ {
			ta[i] = 0
			tb[i] = 0
		}
		uCount = 0
		dfs(-1, -1, 0)
		hh := make([]int, m)
		for h := 0; h < m; h++ {
			if ta[ii[h]] > ta[jj[h]] {
				ii[h], jj[h] = jj[h], ii[h]
			}
			hh[h] = h
		}
		sortEdges(hh, 0, m, compareComp)
		// clear adjacency
		for i := 0; i < n; i++ {
			eh[i] = eh[i][:0]
			inArr[i] = false
		}
		yes := true
		// process components
		i := 0
		for i < m {
			root := findSet(hh[i])
			j := i + 1
			for j < m && findSet(hh[j]) == root {
				j++
			}
			if !solve(hh[i:j], j-i) {
				yes = false
				break
			}
			i = j
		}
		if yes {
			fmt.Fprintln(writer, "YES")
			for u := 0; u < n; u++ {
				for _, h := range eh[u] {
					fmt.Fprintf(writer, "%d ", u^ii[h]^jj[h])
				}
				fmt.Fprintln(writer)
			}
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
