package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader *bufio.Reader

func nextInt() int {
	sign := 1
	val := 0
	c, err := reader.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = reader.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = reader.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = reader.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

type segTree struct {
	size int
	max  []int
	lazy []int
}

func newSegTree(arr []int) *segTree {
	n := len(arr)
	st := &segTree{
		size: n,
		max:  make([]int, n*4),
		lazy: make([]int, n*4),
	}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segTree) build(node, l, r int, arr []int) {
	if l == r {
		st.max[node] = arr[l]
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, arr)
	st.build(node<<1|1, mid+1, r, arr)
	st.pull(node)
}

func (st *segTree) pull(node int) {
	if st.max[node<<1] > st.max[node<<1|1] {
		st.max[node] = st.max[node<<1]
	} else {
		st.max[node] = st.max[node<<1|1]
	}
}

func (st *segTree) apply(node, val int) {
	st.max[node] += val
	st.lazy[node] += val
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		st.apply(node<<1, st.lazy[node])
		st.apply(node<<1|1, st.lazy[node])
		st.lazy[node] = 0
	}
}

func (st *segTree) addRange(l, r, val int) {
	if l > r {
		return
	}
	st.add(1, 0, st.size-1, l, r, val)
}

func (st *segTree) add(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(node, val)
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.add(node<<1, l, mid, ql, qr, val)
	st.add(node<<1|1, mid+1, r, ql, qr, val)
	st.pull(node)
}

func (st *segTree) rightmostNonNegative() int {
	return st.rightmost(1, 0, st.size-1)
}

func (st *segTree) rightmost(node, l, r int) int {
	if st.max[node] < 0 {
		return -1
	}
	if l == r {
		return l
	}
	st.push(node)
	mid := (l + r) >> 1
	if res := st.rightmost(node<<1|1, mid+1, r); res != -1 {
		return res
	}
	return st.rightmost(node<<1, l, mid)
}

func calcS(a, d, ac, dr int) int {
	res := 0
	if a > ac {
		res += a - ac
	}
	if d > dr {
		res += d - dr
	}
	return res
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	ac := nextInt()
	dr := nextInt()
	n := nextInt()

	a := make([]int, n)
	d := make([]int, n)

	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}
	for i := 0; i < n; i++ {
		d[i] = nextInt()
	}

	freq := make([]int, n+1)
	for i := 0; i < n; i++ {
		s := calcS(a[i], d[i], ac, dr)
		if s <= n {
			freq[s]++
		}
	}

	diff := make([]int, n+1)
	prefix := 0
	for i := 0; i <= n; i++ {
		prefix += freq[i]
		diff[i] = prefix - i
	}

	seg := newSegTree(diff)

	m := nextInt()
	for i := 0; i < m; i++ {
		k := nextInt() - 1
		na := nextInt()
		nd := nextInt()

		oldS := calcS(a[k], d[k], ac, dr)
		if oldS <= n {
			seg.addRange(oldS, n, -1)
		}

		a[k] = na
		d[k] = nd

		newS := calcS(na, nd, ac, dr)
		if newS <= n {
			seg.addRange(newS, n, 1)
		}

		ans := seg.rightmostNonNegative()
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
