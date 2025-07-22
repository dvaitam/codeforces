package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Fenwick tree supporting prefix sums of int64 values.
type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) Add(i int, v int64) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int64 {
	if i > f.n {
		i = f.n
	}
	s := int64(0)
	for i > 0 {
		s += f.bit[i]
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
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	// coordinate compression
	vals := make([]int64, n)
	copy(vals, a)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	m := 0
	for i := 0; i < n; i++ {
		if i == 0 || vals[i] != vals[i-1] {
			vals[m] = vals[i]
			m++
		}
	}
	vals = vals[:m]

	fcnt := NewFenwick(m)
	fsum := NewFenwick(m)

	var totalCnt int64
	var totalSum int64
	var ans int64

	for _, x := range a {
		// counts and sums for ai <= x-2
		idx1 := sort.Search(m, func(i int) bool { return vals[i] > x-2 })
		c1 := fcnt.Sum(idx1)
		s1 := fsum.Sum(idx1)
		// counts and sums for ai >= x+2
		idx2 := sort.Search(m, func(i int) bool { return vals[i] > x+1 })
		cprefix := fcnt.Sum(idx2)
		sprefix := fsum.Sum(idx2)
		c2 := totalCnt - cprefix
		s2 := totalSum - sprefix

		ans += int64(x)*(c1+c2) - (s1 + s2)

		// update trees
		pos := sort.Search(m, func(i int) bool { return vals[i] >= x }) + 1
		fcnt.Add(pos, 1)
		fsum.Add(pos, int64(x))
		totalCnt++
		totalSum += int64(x)
	}

	fmt.Fprint(writer, ans)
}
