package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 2e18

type Segment struct {
	id int
	L  int
	R  int
	k  int64
}

type Node struct {
	max_W    int64
	lazy_f   int64
	min_f    int64
	max_f    int64
	max_val  int64
	has_lazy bool
}

var tree []Node

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func pushup(node int) {
	tree[node].max_W = max(tree[node*2].max_W, tree[node*2+1].max_W)
	tree[node].min_f = min(tree[node*2].min_f, tree[node*2+1].min_f)
	tree[node].max_f = max(tree[node*2].max_f, tree[node*2+1].max_f)
	tree[node].max_val = max(tree[node*2].max_val, tree[node*2+1].max_val)
}

func apply_lazy(node int, V int64) {
	tree[node].lazy_f = V
	tree[node].has_lazy = true
	tree[node].min_f = V
	tree[node].max_f = V
	if tree[node].max_W == -INF {
		tree[node].max_val = -INF
	} else {
		tree[node].max_val = V + tree[node].max_W
	}
}

func pushdown(node int) {
	if tree[node].has_lazy {
		apply_lazy(node*2, tree[node].lazy_f)
		apply_lazy(node*2+1, tree[node].lazy_f)
		tree[node].has_lazy = false
	}
}

func build(node, L, R int) {
	tree[node].max_W = -INF
	tree[node].min_f = -INF
	tree[node].max_f = -INF
	tree[node].max_val = -INF
	tree[node].lazy_f = -INF
	tree[node].has_lazy = false
	if L == R {
		return
	}
	mid := (L + R) / 2
	build(node*2, L, mid)
	build(node*2+1, mid+1, R)
}

func update_W(node, L, R, pos int, W int64) {
	if L == R {
		tree[node].max_W = W
		if W == -INF {
			tree[node].max_val = -INF
		} else {
			if tree[node].min_f == -INF {
				tree[node].max_val = -INF
			} else {
				tree[node].max_val = tree[node].min_f + W
			}
		}
		return
	}
	pushdown(node)
	mid := (L + R) / 2
	if pos <= mid {
		update_W(node*2, L, mid, pos, W)
	} else {
		update_W(node*2+1, mid+1, R, pos, W)
	}
	pushup(node)
}

func chmax(node, L, R, ql, qr int, V int64) {
	if ql <= L && R <= qr {
		if tree[node].min_f >= V {
			return
		}
		if tree[node].max_f <= V {
			apply_lazy(node, V)
			return
		}
	}
	pushdown(node)
	mid := (L + R) / 2
	if ql <= mid {
		chmax(node*2, L, mid, ql, qr, V)
	}
	if qr > mid {
		chmax(node*2+1, mid+1, R, ql, qr, V)
	}
	pushup(node)
}

func query_max_val(node, L, R, ql, qr int) int64 {
	if ql <= L && R <= qr {
		return tree[node].max_val
	}
	pushdown(node)
	mid := (L + R) / 2
	res := -INF
	if ql <= mid {
		res = max(res, query_max_val(node*2, L, mid, ql, qr))
	}
	if qr > mid {
		res = max(res, query_max_val(node*2+1, mid+1, R, ql, qr))
	}
	return res
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024*8)
	var n, q int
	fmt.Fscan(reader, &n, &q)

	P1 := make([]int64, n+2)
	P2 := make([]int64, n+2)
	P3 := make([]int64, n+2)

	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P1[i] = P1[i-1] + val
	}
	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P2[i] = P2[i-1] + val
	}
	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P3[i] = P3[i-1] + val
	}

	X := make([]int64, n+1)
	for i := 0; i <= n-1; i++ {
		X[i] = P1[i+1] - P2[i]
	}
	Y := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		Y[i] = P2[i] - P3[i-1]
	}

	segments := make([]Segment, q)
	add := make([][]Segment, n+2)
	exp := make([][]int, n+2)

	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &segments[i].L, &segments[i].R, &segments[i].k)
		segments[i].id = i
		add[segments[i].L] = append(add[segments[i].L], segments[i])
		exp[segments[i].R] = append(exp[segments[i].R], i)
	}

	for i := 1; i <= n; i++ {
		sort.Slice(add[i], func(a, b int) bool {
			return add[i][a].k < add[i][b].k
		})
	}

	ptr := make([]int, n+2)
	deleted := make([]bool, q)

	get_max_W := func(L int) int64 {
		for ptr[L] < len(add[L]) && deleted[add[L][ptr[L]].id] {
			ptr[L]++
		}
		if ptr[L] < len(add[L]) {
			return -add[L][ptr[L]].k
		}
		return -INF
	}

	tree = make([]Node, 4*n+4)
	build(1, 1, n)

	dp := make([]int64, n+2)
	for i := range dp {
		dp[i] = -INF
	}

	ans := -INF

	for i := 1; i <= n; i++ {
		for _, id := range exp[i-1] {
			deleted[id] = true
			L := segments[id].L
			update_W(1, 1, n, L, get_max_W(L))
		}

		V := X[i-1]
		if dp[i-1] > V {
			V = dp[i-1]
		}
		chmax(1, 1, n, 1, i, V)

		update_W(1, 1, n, i, get_max_W(i))

		dp[i] = query_max_val(1, 1, n, 1, i)
		if dp[i] != -INF {
			ans = max(ans, dp[i]+Y[i])
		}
	}

	ans += P3[n]
	fmt.Println(ans)
}
