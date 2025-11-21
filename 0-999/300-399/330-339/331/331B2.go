package main

import (
	"bufio"
	"fmt"
	"os"
)

type BIT struct {
	n int
	f []int
}

func newBIT(n int) *BIT {
	return &BIT{n: n, f: make([]int, n+2)}
}

func (b *BIT) add(idx, delta int) {
	for idx <= b.n {
		b.f[idx] += delta
		idx += idx & -idx
	}
}

func (b *BIT) sum(idx int) int {
	if idx > b.n {
		idx = b.n
	}
	res := 0
	for idx > 0 {
		res += b.f[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *BIT) set(idx, val int) {
	cur := b.sum(idx) - b.sum(idx-1)
	diff := val - cur
	if diff != 0 {
		b.add(idx, diff)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		pos[a[i]] = i
	}

	bit := newBIT(n)
	updateBreak := func(val int) {
		if val < 1 || val >= n {
			return
		}
		if pos[val] > pos[val+1] {
			bit.set(val, 1)
		} else {
			bit.set(val, 0)
		}
	}
	for i := 1; i <= n-1; i++ {
		if pos[i] > pos[i+1] {
			bit.set(i, 1)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x == y {
				fmt.Fprintln(out, 1)
				continue
			}
			if x > y {
				x, y = y, x
			}
			sum := bit.sum(y-1) - bit.sum(x-1)
			fmt.Fprintln(out, sum+1)
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x == y {
				continue
			}
			ax, ay := a[x], a[y]
			a[x], a[y] = a[y], a[x]
			pos[ax], pos[ay] = pos[ay], pos[ax]
			values := []int{ax - 1, ax, ay - 1, ay}
			seen := make(map[int]struct{})
			for _, v := range values {
				if _, ok := seen[v]; ok {
					continue
				}
				seen[v] = struct{}{}
				updateBreak(v)
			}
		}
	}
}
