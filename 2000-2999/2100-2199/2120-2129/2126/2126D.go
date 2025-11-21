package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type casino struct {
	l, r    int
	realIdx int
}

type segTree struct {
	tree [][]int
	n    int
}

func newSegTree(n int) *segTree {
	size := 4 * n
	tree := make([][]int, size)
	return &segTree{tree: tree, n: n}
}

func (st *segTree) add(node, nl, nr, l, r, id int) {
	if l > nr || r < nl {
		return
	}
	if l <= nl && nr <= r {
		st.tree[node] = append(st.tree[node], id)
		return
	}
	mid := (nl + nr) >> 1
	st.add(node<<1, nl, mid, l, r, id)
	st.add(node<<1|1, mid+1, nr, l, r, id)
}

func (st *segTree) query(node, nl, nr, pos int, handler func(int)) {
	for len(st.tree[node]) > 0 {
		id := st.tree[node][len(st.tree[node])-1]
		st.tree[node] = st.tree[node][:len(st.tree[node])-1]
		handler(id)
	}
	if nl == nr {
		return
	}
	mid := (nl + nr) >> 1
	if pos <= mid {
		st.query(node<<1, nl, mid, pos, handler)
	} else {
		st.query(node<<1|1, mid+1, nr, pos, handler)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int
		fmt.Fscan(in, &n, &k)
		casinos := make([]casino, 0, n)
		reals := make([]int, 0, n+1)
		reals = append(reals, k)
		tmpL := make([]int, n)
		tmpR := make([]int, n)
		tmpReal := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tmpL[i], &tmpR[i], &tmpReal[i])
			reals = append(reals, tmpReal[i])
		}
		sort.Ints(reals)
		reals = unique(reals)
		idxOf := make(map[int]int, len(reals))
		for i, v := range reals {
			idxOf[v] = i
		}
		for i := 0; i < n; i++ {
			l := tmpL[i]
			r := tmpR[i]
			realVal := tmpReal[i]
			left := lowerBound(reals, l)
			right := upperBound(reals, r) - 1
			if left <= right {
				casinos = append(casinos, casino{l: left, r: right, realIdx: idxOf[realVal]})
			}
		}
		m := len(reals)
		st := newSegTree(m)
		for id, c := range casinos {
			st.add(1, 0, m-1, c.l, c.r, id)
		}
		visitedVals := make([]bool, m)
		visitedVals[idxOf[k]] = true
		queue := []int{idxOf[k]}
		usedCasino := make([]bool, len(casinos))
		handler := func(id int) {
			if usedCasino[id] {
				return
			}
			usedCasino[id] = true
			realIdx := casinos[id].realIdx
			if !visitedVals[realIdx] {
				visitedVals[realIdx] = true
				queue = append(queue, realIdx)
			}
		}
		head := 0
		for head < len(queue) {
			pos := queue[head]
			head++
			st.query(1, 0, m-1, pos, handler)
		}
		ans := k
		for i, v := range reals {
			if visitedVals[i] && v > ans {
				ans = v
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func lowerBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func upperBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}
