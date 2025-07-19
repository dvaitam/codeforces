package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Item represents an update or query event
type Item struct {
	t, coef, time, pos, val int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		a[x] = i
	}
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		b[i] = a[x]
	}
	// Prepare events
	q := make([]Item, 0, n+4*m)
	for i := 1; i <= n; i++ {
		q = append(q, Item{t: 0, coef: 1, time: 0, pos: i, val: b[i]})
	}
	ans := make([]int64, m+1)
	sentinel := int64(-200233)
	for i := 1; i <= m; i++ {
		var op int
		fmt.Fscan(reader, &op)
		if op == 1 {
			var vl, vr, pl, pr int
			fmt.Fscan(reader, &vl, &vr, &pl, &pr)
			q = append(q, Item{t: i, coef: 1, time: i, pos: pr, val: vr})
			if pl > 1 {
				q = append(q, Item{t: i, coef: -1, time: i, pos: pl - 1, val: vr})
			}
			if vl > 1 {
				q = append(q, Item{t: i, coef: -1, time: i, pos: pr, val: vl - 1})
			}
			if pl > 1 && vl > 1 {
				q = append(q, Item{t: i, coef: 1, time: i, pos: pl - 1, val: vl - 1})
			}
		} else {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			q = append(q, Item{t: 0, coef: 1, time: i, pos: x, val: b[y]})
			q = append(q, Item{t: 0, coef: 1, time: i, pos: y, val: b[x]})
			q = append(q, Item{t: 0, coef: -1, time: i, pos: x, val: b[x]})
			q = append(q, Item{t: 0, coef: -1, time: i, pos: y, val: b[y]})
			b[x], b[y] = b[y], b[x]
			ans[i] = sentinel
		}
	}
	// Sort events by time
	sort.Slice(q, func(i, j int) bool {
		return q[i].time < q[j].time
	})
	// Fenwick tree
	C := make([]int, n+1)
	add := func(p, x int) {
		for i := p; i <= n; i += i & -i {
			C[i] += x
		}
	}
	sum := func(p int) int {
		s := 0
		for i := p; i > 0; i -= i & -i {
			s += C[i]
		}
		return s
	}
	// CDQ divide and conquer
	aux := make([]Item, len(q))
	var cdq func(l, r int)
	cdq = func(l, r int) {
		if l >= r {
			return
		}
		mid := (l + r) >> 1
		cdq(l, mid)
		cdq(mid+1, r)
		i, j, k := l, mid+1, l
		for i <= mid && j <= r {
			if q[i].pos <= q[j].pos {
				if q[i].t == 0 {
					add(q[i].val, q[i].coef)
				}
				aux[k] = q[i]
				i++
			} else {
				if q[j].t != 0 {
					ans[q[j].t] += int64(q[j].coef) * int64(sum(q[j].val))
				}
				aux[k] = q[j]
				j++
			}
			k++
		}
		for i <= mid {
			if q[i].t == 0 {
				add(q[i].val, q[i].coef)
			}
			aux[k] = q[i]
			i++
			k++
		}
		for j <= r {
			if q[j].t != 0 {
				ans[q[j].t] += int64(q[j].coef) * int64(sum(q[j].val))
			}
			aux[k] = q[j]
			j++
			k++
		}
		// rollback
		for p := l; p <= mid; p++ {
			if q[p].t == 0 {
				add(q[p].val, -q[p].coef)
			}
		}
		// restore
		for p := l; p <= r; p++ {
			q[p] = aux[p]
		}
	}
	cdq(0, len(q)-1)
	// Output answers for type-1 queries
	for i := 1; i <= m; i++ {
		if ans[i] != sentinel {
			fmt.Fprintln(writer, ans[i])
		}
	}
}
