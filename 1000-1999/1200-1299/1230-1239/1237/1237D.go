package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf = int(1e9)

// Fenwick tree supporting prefix minimum queries
type fenwick struct {
	n    int
	data []int
}

func newFenwick(n int) *fenwick {
	f := &fenwick{n: n, data: make([]int, n+2)}
	for i := range f.data {
		f.data[i] = inf
	}
	return f
}

func (f *fenwick) update(i, val int) {
	i++
	for i <= f.n+1 {
		if val < f.data[i] {
			f.data[i] = val
		}
		i += i & -i
	}
}

func (f *fenwick) query(i int) int {
	i++
	res := inf
	for i > 0 {
		if f.data[i] < res {
			res = f.data[i]
		}
		i -= i & -i
	}
	return res
}

func uniqueSorted(a []int) []int {
	if len(a) == 0 {
		return a
	}
	sort.Ints(a)
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	minv, maxv := inf, 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] < minv {
			minv = a[i]
		}
		if a[i] > maxv {
			maxv = a[i]
		}
	}
	if minv*2 >= maxv {
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, -1)
		}
		fmt.Fprintln(out)
		return
	}

	m := 3 * n
	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = a[i%n]
	}

	vals := make([]int, m)
	copy(vals, b)
	vals = uniqueSorted(vals)
	idx := make(map[int]int, len(vals))
	for i, v := range vals {
		idx[v] = i
	}
	bit := newFenwick(len(vals))

	nextSmall := make([]int, m)
	for i := m - 1; i >= 0; i-- {
		thr := (b[i] - 1) / 2
		pos := sort.SearchInts(vals, thr+1) - 1
		if pos >= 0 {
			nextSmall[i] = bit.query(pos)
		} else {
			nextSmall[i] = inf
		}
		bit.update(idx[b[i]], i)
	}

	nextGreater := make([]int, m)
	stack := make([]int, 0)
	for i := m - 1; i >= 0; i-- {
		for len(stack) > 0 && b[stack[len(stack)-1]] <= b[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			nextGreater[i] = stack[len(stack)-1]
		} else {
			nextGreater[i] = inf
		}
		stack = append(stack, i)
	}

	ans := make([]int, m)
	for i := range ans {
		ans[i] = inf
	}
	for i := m - 1; i >= 0; i-- {
		ns, ng := nextSmall[i], nextGreater[i]
		if ns < ng {
			ans[i] = ns - i
		} else {
			if ng >= m || ans[ng] >= inf {
				ans[i] = inf
			} else {
				ans[i] = ans[ng] + (ng - i)
			}
		}
	}

	for i := 0; i < n; i++ {
		if ans[i] >= inf {
			ans[i] = -1
		}
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
