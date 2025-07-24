package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	var s int64
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func (b *BIT) LowerBound(target int64) int {
	idx := 0
	bitMask := 1 << (bits.Len(uint(b.n)) - 1)
	for bitMask > 0 {
		next := idx + bitMask
		if next <= b.n && b.tree[next] <= target {
			target -= b.tree[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	s := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &s[i])
	}

	bit := NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(i, int64(i))
	}

	p := make([]int, n+1)
	for i := n; i >= 1; i-- {
		idx := bit.LowerBound(s[i])
		p[i] = idx
		bit.Add(idx, -int64(idx))
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, p[i])
	}
	writer.WriteByte('\n')
}
