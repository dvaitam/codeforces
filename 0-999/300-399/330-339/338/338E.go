package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Scanner struct {
	b   []byte
	pos int
}

func NewScanner() *Scanner {
	b, _ := io.ReadAll(os.Stdin)
	return &Scanner{b: b, pos: 0}
}

func (s *Scanner) nextInt() int {
	for s.pos < len(s.b) && (s.b[s.pos] < '0' || s.b[s.pos] > '9') {
		s.pos++
	}
	if s.pos >= len(s.b) {
		return 0
	}
	res := 0
	for s.pos < len(s.b) && s.b[s.pos] >= '0' && s.b[s.pos] <= '9' {
		res = res*10 + int(s.b[s.pos]-'0')
		s.pos++
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var tree []int
var lazy []int

func build(node, l, r int) {
	if l == r {
		tree[node] = -l
		return
	}
	mid := (l + r) / 2
	build(node*2, l, mid)
	build(node*2+1, mid+1, r)
	tree[node] = min(tree[node*2], tree[node*2+1])
}

func push(node int) {
	if lazy[node] != 0 {
		tree[node*2] += lazy[node]
		lazy[node*2] += lazy[node]
		tree[node*2+1] += lazy[node]
		lazy[node*2+1] += lazy[node]
		lazy[node] = 0
	}
}

func update(node, l, r, ql, qr, val int) {
	if ql <= l && r <= qr {
		tree[node] += val
		lazy[node] += val
		return
	}
	push(node)
	mid := (l + r) / 2
	if ql <= mid {
		update(node*2, l, mid, ql, qr, val)
	}
	if qr > mid {
		update(node*2+1, mid+1, r, ql, qr, val)
	}
	tree[node] = min(tree[node*2], tree[node*2+1])
}

func getIdx(x int, c []int, m int) int {
	left, right := 1, m
	ans := m + 1
	for left <= right {
		mid := (left + right) / 2
		if c[mid] <= x {
			ans = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return ans
}

func main() {
	sc := NewScanner()
	n := sc.nextInt()
	if n == 0 {
		return
	}
	m := sc.nextInt()
	h := sc.nextInt()

	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = sc.nextInt()
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = sc.nextInt()
	}

	sort.Ints(b)
	c := make([]int, m+1)
	for i := 1; i <= m; i++ {
		c[i] = h - b[i-1]
	}

	tree = make([]int, 4*m+5)
	lazy = make([]int, 4*m+5)

	build(1, 1, m)

	ans := 0
	for i := 0; i < m; i++ {
		idx := getIdx(a[i], c, m)
		if idx <= m {
			update(1, 1, m, idx, m, 1)
		}
	}

	if tree[1] >= 0 {
		ans++
	}

	for i := m; i < n; i++ {
		idxOut := getIdx(a[i-m], c, m)
		if idxOut <= m {
			update(1, 1, m, idxOut, m, -1)
		}

		idxIn := getIdx(a[i], c, m)
		if idxIn <= m {
			update(1, 1, m, idxIn, m, 1)
		}

		if tree[1] >= 0 {
			ans++
		}
	}

	fmt.Println(ans)
}
