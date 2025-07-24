package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int = 1000000007

type Barrier struct {
	u int
	l int
	r int
	s int
}

type pair struct {
	u   int
	idx int
}

var (
	w        int
	barriers []Barrier
	tree     [][]pair
	memo     map[[2]int]int
)

func add(node, L, R, l, r int, p pair) {
	if l <= L && R <= r {
		tree[node] = append(tree[node], p)
		return
	}
	mid := (L + R) / 2
	if l <= mid {
		add(node*2, L, mid, l, r, p)
	}
	if r > mid {
		add(node*2+1, mid+1, R, l, r, p)
	}
}

func build() {
	tree = make([][]pair, 4*(w+2))
	for i, b := range barriers {
		add(1, 1, w, b.l, b.r, pair{b.u, i})
	}
	for i := range tree {
		if len(tree[i]) > 1 {
			sort.Slice(tree[i], func(a, b int) bool {
				return tree[i][a].u < tree[i][b].u
			})
		}
	}
}

func queryRec(node, L, R, col, h int, bestIdx *int, bestRow *int) {
	arr := tree[node]
	if len(arr) > 0 {
		pos := sort.Search(len(arr), func(i int) bool { return arr[i].u >= h })
		if pos > 0 {
			cand := arr[pos-1]
			if cand.u > *bestRow {
				*bestRow = cand.u
				*bestIdx = cand.idx
			}
		}
	}
	if L == R {
		return
	}
	mid := (L + R) / 2
	if col <= mid {
		queryRec(node*2, L, mid, col, h, bestIdx, bestRow)
	} else {
		queryRec(node*2+1, mid+1, R, col, h, bestIdx, bestRow)
	}
}

func nextBarrier(col, h int) int {
	bestIdx := -1
	bestRow := -1
	queryRec(1, 1, w, col, h, &bestIdx, &bestRow)
	return bestIdx
}

func solve(col, h int) int {
	key := [2]int{col, h}
	if val, ok := memo[key]; ok {
		return val
	}
	cur := h
	for {
		idx := nextBarrier(col, cur)
		if idx == -1 {
			memo[key] = 1
			return 1
		}
		b := barriers[idx]
		if h > b.u+b.s {
			cur = b.u
			continue
		}
		left := b.l - 1
		right := b.r + 1
		if left < 1 {
			left = right
		}
		if right > w {
			right = left
		}
		var res int
		if left == right {
			res = solve(left, b.u)
			res = (res * 2) % mod
		} else {
			res = solve(left, b.u)
			res += solve(right, b.u)
			if res >= mod {
				res -= mod
			}
		}
		memo[key] = res
		return res
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var h, n int
	fmt.Fscan(reader, &h, &w, &n)
	barriers = make([]Barrier, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &barriers[i].u, &barriers[i].l, &barriers[i].r, &barriers[i].s)
	}

	build()
	memo = make(map[[2]int]int)
	start := h + 1
	ans := 0
	for c := 1; c <= w; c++ {
		ans += solve(c, start)
		if ans >= mod {
			ans -= mod
		}
	}
	fmt.Fprintln(writer, ans%mod)
}
