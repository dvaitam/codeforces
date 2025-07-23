package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	n int
	t []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, t: make([]int64, n+2)}
}

func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.t[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += b.t[i]
		i -= i & -i
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > n {
			a[i] = n
		}
	}

	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i + 1
	}
	sort.Slice(idx, func(i, j int) bool { return a[idx[i]] > a[idx[j]] })

	bit := NewBIT(n)
	ptr := 0
	ans := int64(0)

	for j := n; j >= 1; j-- {
		for ptr < n && a[idx[ptr]] >= j {
			bit.Add(idx[ptr], 1)
			ptr++
		}
		upper := a[j]
		if upper > j-1 {
			upper = j - 1
		}
		if upper > 0 {
			ans += bit.Sum(upper)
		}
	}

	fmt.Fprintln(writer, ans)
}
