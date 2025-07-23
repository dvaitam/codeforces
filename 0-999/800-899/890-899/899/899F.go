package main

import (
	"bufio"
	"fmt"
	"os"
)

// BIT implements a Fenwick tree for prefix sums and order statistics.
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, delta int) {
	for x := i; x <= b.n; x += x & -x {
		b.tree[x] += delta
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}

// FindKth returns smallest index i such that Sum(i) >= k.
func (b *BIT) FindKth(k int) int {
	pos := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		nxt := pos + d
		if nxt <= b.n && b.tree[nxt] < k {
			k -= b.tree[nxt]
			pos = nxt
		}
	}
	return pos + 1
}

// FindNextFrom returns the first index >= l having value 1; if none, returns n+1.
func (b *BIT) FindNextFrom(l int) int {
	pre := b.Sum(l - 1)
	total := b.Sum(b.n)
	if pre == total {
		return b.n + 1
	}
	return b.FindKth(pre + 1)
}

func charIndex(ch byte) int {
	if ch >= 'a' && ch <= 'z' {
		return int(ch - 'a')
	}
	if ch >= 'A' && ch <= 'Z' {
		return 26 + int(ch-'A')
	}
	return 52 + int(ch-'0')
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	bits := make([]*BIT, 62)
	for i := 0; i < 62; i++ {
		bits[i] = NewBIT(n)
	}
	global := NewBIT(n)
	alive := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		global.Add(i, 1)
		idx := charIndex(s[i-1])
		bits[idx].Add(i, 1)
		alive[i] = true
	}

	for op := 0; op < m; op++ {
		var l, r int
		var cs string
		fmt.Fscan(reader, &l, &r, &cs)
		c := cs[0]
		idx := charIndex(c)
		lIdx := global.FindKth(l)
		rIdx := global.FindKth(r)
		limit := rIdx
		pos := bits[idx].FindNextFrom(lIdx)
		for pos <= limit {
			if alive[pos] {
				alive[pos] = false
				global.Add(pos, -1)
				bits[idx].Add(pos, -1)
			}
			pos = bits[idx].FindNextFrom(pos + 1)
		}
	}

	res := make([]byte, 0, n)
	for i := 1; i <= n; i++ {
		if alive[i] {
			res = append(res, s[i-1])
		}
	}
	writer.Write(res)
}
