package main

import (
	"bufio"
	"fmt"
	"os"
)

// BIT implements a Fenwick tree (binary indexed tree).
type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) Update(i int, delta int64) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Query(i int) int64 {
	var res int64
	for i > 0 {
		res += b.tree[i]
		i -= i & -i
	}
	return res
}

// Kth returns the smallest index x such that sum_{i<=x} >= k.
func (b *BIT) Kth(k int64) int {
	idx := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	pos := make([]int, n+1)
	for i, v := range p {
		pos[v] = i + 1
	}

	bitCount := NewBIT(n)
	bitSum := NewBIT(n)

	results := make([]int64, n)
	var inversions int64
	var sumTotal int64

	for k := 1; k <= n; k++ {
		v := pos[k]
		inversions += int64(k-1) - bitCount.Query(v)
		bitCount.Update(v, 1)
		bitSum.Update(v, int64(v))
		sumTotal += int64(v)

		m := (k + 1) / 2
		rm := bitCount.Kth(int64(m))
		sumFirstM := bitSum.Query(rm)

		left := int64(m*(m-1)) / 2
		right := int64((k-m)*(k-m+1)) / 2
		costGroup := sumTotal - 2*sumFirstM + int64(2*m-k)*int64(rm) - left - right

		results[k-1] = costGroup + inversions
	}

	out := bufio.NewWriter(os.Stdout)
	for i, v := range results {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
	out.Flush()
}
