package main

import (
	"bufio"
	"fmt"
	"os"
)

type fenwick struct {
	n int
	t []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, t: make([]int, n+2)}
}

func (f *fenwick) add(idx, val int) {
	for idx <= f.n {
		f.t[idx] += val
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.t[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		fw := newFenwick(n + 2)
		prevGreater := make([]int, n)
		var base int64
		for i := 0; i < n; i++ {
			val := p[i]
			sm := fw.sum(val)
			prevGreater[i] = i - sm
			base += int64(prevGreater[i])
			fw.add(val, 1)
		}
		fw = newFenwick(n + 2)
		nextGreater := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			val := p[i]
			nextGreater[i] = fw.sum(n+1) - fw.sum(val)
			fw.add(val, 1)
		}
		ans := base
		for i := 0; i < n; i++ {
			diff := nextGreater[i] - prevGreater[i]
			if diff < 0 {
				ans += int64(diff)
			}
		}
		fmt.Fprintf(out, "%d \n", ans)
	}
}
