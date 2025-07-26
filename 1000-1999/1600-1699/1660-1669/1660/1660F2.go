package main

import (
	"bufio"
	"fmt"
	"os"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n + 2, tree: make([]int, n+3)}
	return b
}

func (b *BIT) Add(idx, val int) {
	idx++
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	if idx < 0 {
		return 0
	}
	if idx >= b.n {
		idx = b.n - 1
	}
	idx++
	res := 0
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		size := 2*n + 5
		bits := []*BIT{NewBIT(size), NewBIT(size), NewBIT(size)}
		offset := n + 2
		sum := 0
		bits[sum%3].Add(sum+offset, 1)
		var ans int64
		for i := 0; i < n; i++ {
			if s[i] == '+' {
				sum--
			} else {
				sum++
			}
			mod := ((sum % 3) + 3) % 3
			idx := sum + offset
			ans += int64(bits[mod].Sum(idx))
			bits[mod].Add(idx, 1)
		}
		fmt.Fprintln(writer, ans)
	}
}
