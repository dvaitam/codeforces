package main

import (
	"bufio"
	"fmt"
	"os"
)

type bit struct {
	n    int
	tree []int
}

func newBIT(n int) *bit {
	return &bit{n: n, tree: make([]int, n+2)}
}

func (b *bit) add(idx, val int) {
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *bit) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b.tree[idx]
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
		pref := make([]int64, n+1)
		b1 := newBIT(n + 2)
		for i := 1; i <= n; i++ {
			val := p[i-1]
			greater := (i - 1) - b1.sum(val)
			pref[i] = pref[i-1] + int64(greater)
			b1.add(val, 1)
		}
		suff := make([]int64, n+2)
		b2 := newBIT(2*n + 2)
		for i := n - 1; i >= 0; i-- {
			val := 2*n - p[i]
			less := b2.sum(val - 1)
			suff[i+1] = suff[i+2] + int64(less)
			b2.add(val, 1)
		}
		ans := pref[n] + suff[n+1]
		for t := 0; t <= n; t++ {
			val := pref[t] + suff[t+1]
			if val < ans {
				ans = val
			}
		}
		fmt.Fprintf(out, "%d \n", ans)
	}
}
