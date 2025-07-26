package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// Fenwick tree for prefix sums and searching
type Fenwick struct {
	n int
	t []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, t: make([]int, n+1)}
}

func (f *Fenwick) Add(idx, val int) {
	for i := idx + 1; i <= f.n; i += i & -i {
		f.t[i] += val
	}
}

func (f *Fenwick) Sum(idx int) int {
	if idx < 0 {
		return 0
	}
	res := 0
	for i := idx + 1; i > 0; i -= i & -i {
		res += f.t[i]
	}
	return res
}

func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func (f *Fenwick) FindFirst() int {
	total := f.Sum(f.n - 1)
	if total == 0 {
		return -1
	}
	idx := 0
	bit := 1 << (bits.Len(uint(f.n)))
	sum := 0
	for bit > 0 {
		next := idx + bit
		if next <= f.n && sum+f.t[next] == 0 {
			idx = next
			sum += f.t[next]
		}
		bit >>= 1
	}
	return idx
}

func (f *Fenwick) FindLast() int {
	total := f.Sum(f.n - 1)
	if total == 0 {
		return -1
	}
	idx := 0
	bit := 1 << (bits.Len(uint(f.n)))
	sum := 0
	target := total - 1
	for bit > 0 {
		next := idx + bit
		if next <= f.n && sum+f.t[next] <= target {
			idx = next
			sum += f.t[next]
		}
		bit >>= 1
	}
	return idx
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Segment tree storing sum and minimum prefix
type SegTree struct {
	n   int
	sum []int
	pre []int
}

func NewSegTree(arr []int) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &SegTree{n: n, sum: make([]int, 2*n), pre: make([]int, 2*n)}
	for i := 0; i < len(arr); i++ {
		st.sum[n+i] = arr[i]
		if arr[i] < 0 {
			st.pre[n+i] = arr[i]
		} else {
			st.pre[n+i] = 0
		}
	}
	for i := n - 1; i > 0; i-- {
		st.pull(i)
	}
	return st
}

func (st *SegTree) pull(v int) {
	l, r := v<<1, v<<1|1
	st.sum[v] = st.sum[l] + st.sum[r]
	st.pre[v] = min(st.pre[l], st.sum[l]+st.pre[r])
}

func (st *SegTree) Update(idx, val int) {
	v := st.n + idx
	st.sum[v] = val
	if val < 0 {
		st.pre[v] = val
	} else {
		st.pre[v] = 0
	}
	for v >>= 1; v > 0; v >>= 1 {
		st.pull(v)
	}
}

func (st *SegTree) query(v, l, r, L, R int) (int, int) {
	if L <= l && r <= R {
		return st.sum[v], st.pre[v]
	}
	m := (l + r) >> 1
	if R <= m {
		return st.query(v<<1, l, m, L, R)
	}
	if L > m {
		return st.query(v<<1|1, m+1, r, L, R)
	}
	sumL, preL := st.query(v<<1, l, m, L, R)
	sumR, preR := st.query(v<<1|1, m+1, r, L, R)
	return sumL + sumR, min(preL, sumL+preR)
}

func (st *SegTree) Prefix(idx int) (int, int) {
	if idx < 0 {
		return 0, 0
	}
	if idx >= st.n {
		idx = st.n - 1
	}
	return st.query(1, 0, st.n-1, 0, idx)
}

func (st *SegTree) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	s, _ := st.query(1, 0, st.n-1, l, r)
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, q int
	fmt.Fscan(in, &n, &q)
	var str string
	fmt.Fscan(in, &str)
	s := []byte(str)

	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if s[i] == '(' {
			arr[i] = 1
		} else {
			arr[i] = -1
		}
	}
	st := NewSegTree(arr)

	pairOpen := make([]bool, n-1)
	pairClose := make([]bool, n-1)
	fenOpen := NewFenwick(n - 1)
	fenClose := NewFenwick(n - 1)
	for i := 0; i < n-1; i++ {
		if s[i] == '(' && s[i+1] == '(' {
			pairOpen[i] = true
			fenOpen.Add(i, 1)
		}
		if s[i] == ')' && s[i+1] == ')' {
			pairClose[i] = true
			fenClose.Add(i, 1)
		}
	}

	check := func() bool {
		if n%2 == 1 || s[0] != '(' || s[n-1] != ')' {
			return false
		}
		total, minPref := st.Prefix(n - 1)
		if total == 0 && minPref >= 0 {
			return true
		}
		first := fenOpen.FindFirst()
		last := fenClose.FindLast()
		if first == -1 || last == -1 || first >= last {
			return false
		}
		_, pref := st.Prefix(first)
		if pref < 0 {
			return false
		}
		if st.RangeSum(last+2, n-1) > 0 {
			return false
		}
		return true
	}

	for ; q > 0; q-- {
		var pos int
		fmt.Fscan(in, &pos)
		pos--
		// flip char
		if s[pos] == '(' {
			s[pos] = ')'
			st.Update(pos, -1)
		} else {
			s[pos] = '('
			st.Update(pos, 1)
		}
		// update pairs around pos
		for i := pos - 1; i <= pos; i++ {
			if i >= 0 && i < n-1 {
				newOpen := s[i] == '(' && s[i+1] == '('
				if pairOpen[i] != newOpen {
					if pairOpen[i] {
						fenOpen.Add(i, -1)
					} else {
						fenOpen.Add(i, 1)
					}
					pairOpen[i] = newOpen
				}
				newClose := s[i] == ')' && s[i+1] == ')'
				if pairClose[i] != newClose {
					if pairClose[i] {
						fenClose.Add(i, -1)
					} else {
						fenClose.Add(i, 1)
					}
					pairClose[i] = newClose
				}
			}
		}
		if check() {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
