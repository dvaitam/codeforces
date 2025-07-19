package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Q represents an event: segment addition (t>0) or query (t<0)
type Q struct {
	t    int   // >0 for segment idx, <0 for query -idx
	x, y int64 // sort keys
	L, R int   // compressed coordinates
}

var (
	q   []Q
	tmp []Q
	bit []int
	ans []int
)

func bitUpdate(i, v int) {
	n := len(bit)
	for ; i < n; i += i & -i {
		bit[i] += v
	}
}

func bitQuery(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += bit[i]
	}
	return s
}

func solveCDQ(l, r int) {
	if l >= r {
		return
	}
	m := (l + r) >> 1
	solveCDQ(l, m)
	solveCDQ(m+1, r)
	i, j := l, m+1
	k := l
	// merge by y, and process
	for i <= m || j <= r {
		if j > r || (i <= m && (q[i].y < q[j].y || (q[i].y == q[j].y && q[i].t > 0))) {
			// left event
			tmp[k] = q[i]
			if q[i].t > 0 {
				bitUpdate(q[i].L, 1)
				bitUpdate(q[i].R+1, -1)
			}
			i++
		} else {
			// right event (query)
			tmp[k] = q[j]
			if q[j].t < 0 {
				idx := -q[j].t
				ans[idx] += bitQuery(q[j].L)
			}
			j++
		}
		k++
	}
	// clear BIT
	for p := l; p <= m; p++ {
		if q[p].t > 0 {
			bitUpdate(q[p].L, -1)
			bitUpdate(q[p].R+1, 1)
		}
	}
	// copy back
	for p := l; p <= r; p++ {
		q[p] = tmp[p]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	a := make([]int64, n)
	b := make([]int64, n)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}
	d := make([]int64, m)
	e := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &d[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &e[i])
	}
	// build events
	total := n + m
	q = make([]Q, 0, total)
	vals := make([]int64, 0, 2*total)
	for i := 0; i < n; i++ {
		x := a[i] - c[i]
		y := a[i] + c[i]
		q = append(q, Q{t: i + 1, x: x, y: y, L: int(a[i]), R: int(b[i])})
		vals = append(vals, a[i], b[i])
	}
	for i := 0; i < m; i++ {
		x := d[i] - e[i]
		y := d[i] + e[i]
		// L=R=d[i]
		q = append(q, Q{t: -(i + 1), x: x, y: y, L: int(d[i]), R: int(d[i])})
		vals = append(vals, d[i], d[i])
	}
	// coordinate compression
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	uni := vals[:0]
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			uni = append(uni, v)
		}
	}
	// map L, R
	for idx := range q {
		// find by binary search on original L, R
		origL := q[idx].L
		origR := q[idx].R
		li := sort.Search(len(uni), func(i int) bool { return uni[i] >= int64(origL) })
		ri := sort.Search(len(uni), func(i int) bool { return uni[i] >= int64(origR) })
		q[idx].L = li + 1
		q[idx].R = ri + 1
	}
	// initial sort by x, y, t desc
	sort.Slice(q, func(i, j int) bool {
		if q[i].x != q[j].x {
			return q[i].x < q[j].x
		}
		if q[i].y != q[j].y {
			return q[i].y < q[j].y
		}
		return q[i].t > q[j].t
	})
	tmp = make([]Q, len(q))
	// initialize BIT and ans
	size := len(uni) + 5
	bit = make([]int, size)
	ans = make([]int, m+1)
	// divide and conquer
	solveCDQ(0, len(q)-1)
	// output
	for i := 1; i <= m; i++ {
		fmt.Fprint(out, ans[i])
		if i < m {
			out.WriteByte(' ')
		}
	}
	out.WriteByte('\n')
}
