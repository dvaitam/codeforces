package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Fenwick tree for prefix sums
type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	if i > f.n {
		i = f.n
	}
	s := 0
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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			vals[i] = arr[i]
		}
		sort.Ints(vals)
		vals = unique(vals)
		m := len(vals)
		comp := make(map[int]int, m)
		for i, v := range vals {
			comp[v] = i + 1 // 1-indexed for Fenwick
		}

		ft := NewFenwick(m)
		total := 0
		ans := 0
		for _, v := range arr {
			idx := comp[v]
			less := ft.Sum(idx - 1)
			greater := total - ft.Sum(idx)
			if less < greater {
				ans += less
			} else {
				ans += greater
			}
			ft.Add(idx, 1)
			total++
		}
		fmt.Fprintln(writer, ans)
	}
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}
