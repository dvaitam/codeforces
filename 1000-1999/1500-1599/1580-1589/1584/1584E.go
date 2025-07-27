package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type SegTree struct {
	n    int
	size int
	tree []int
}

func NewSegTree(n int, inf int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int, 2*size)
	for i := range tree {
		tree[i] = inf
	}
	return &SegTree{n: n, size: size, tree: tree}
}

func (st *SegTree) Update(pos, val int) {
	p := pos + st.size
	if val < st.tree[p] {
		st.tree[p] = val
		p >>= 1
		for p > 0 {
			if st.tree[2*p] < st.tree[2*p+1] {
				st.tree[p] = st.tree[2*p]
			} else {
				st.tree[p] = st.tree[2*p+1]
			}
			p >>= 1
		}
	}
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return st.tree[0] // assume tree[0] = inf??? Wait we not set, we return inf.
	}
	l += st.size
	r += st.size
	res := st.tree[0] // we will set outside
	for l <= r {
		if l&1 == 1 {
			if st.tree[l] < res {
				res = st.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.tree[r] < res {
				res = st.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func solveCase(n int, arr []int) int64 {
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			prefix[i] = prefix[i-1] - int64(arr[i-1])
		} else {
			prefix[i] = prefix[i-1] + int64(arr[i-1])
		}
	}
	vals := make([]int64, n+1)
	copy(vals, prefix)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	vals = uniqueInt64(vals)
	idx := make(map[int64]int, len(vals))
	for i, v := range vals {
		idx[v] = i
	}
	m := len(vals)
	INF := n + 1
	segOdd := NewSegTree(m, INF)
	segEven := NewSegTree(m, INF)
	nextGE := make([]int, n+1)
	nextLE := make([]int, n+1)
	for i := n; i >= 0; i-- {
		r := idx[prefix[i]]
		ng := INF
		if r+1 <= m-1 {
			ng = segOdd.Query(r+1, m-1)
		}
		nextGE[i] = ng
		if i%2 == 1 {
			segOdd.Update(r, i)
		}
		nl := INF
		if r-1 >= 0 {
			nl = segEven.Query(0, r-1)
		}
		nextLE[i] = nl
		if i%2 == 0 {
			segEven.Update(r, i)
		}
	}
	limit := make([]int, n+1)
	for i := 0; i <= n; i++ {
		nxt := nextGE[i]
		if nextLE[i] < nxt {
			nxt = nextLE[i]
		}
		if nxt == INF {
			limit[i] = n
		} else {
			limit[i] = nxt - 1
		}
	}
	pos := make(map[int64][]int)
	for i, v := range prefix {
		pos[v] = append(pos[v], i)
	}
	var ans int64
	for i := 0; i < n; i++ {
		arrPos := pos[prefix[i]]
		l := sort.SearchInts(arrPos, i+1)
		r := sort.SearchInts(arrPos, limit[i]+1)
		ans += int64(r - l)
	}
	return ans
}

func uniqueInt64(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	res := a[:1]
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			res = append(res, a[i])
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		ans := solveCase(n, arr)
		fmt.Fprintln(out, ans)
	}
}
