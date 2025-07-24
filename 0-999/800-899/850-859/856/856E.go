package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Event struct {
	t  int
	x  float64
	y  float64
	i  int
	j  int
	id int
}

type SegTree struct {
	n   int
	min []int
	add []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: size, min: make([]int, size<<1), add: make([]int, size<<1)}
}

func (st *SegTree) push(v int) {
	if st.add[v] != 0 {
		val := st.add[v]
		st.min[v<<1] += val
		st.add[v<<1] += val
		st.min[v<<1|1] += val
		st.add[v<<1|1] += val
		st.add[v] = 0
	}
}

func (st *SegTree) rangeAdd(v, l, r, ql, qr, delta int) {
	if ql <= l && r <= qr {
		st.min[v] += delta
		st.add[v] += delta
		return
	}
	st.push(v)
	m := (l + r) >> 1
	if ql < m {
		st.rangeAdd(v<<1, l, m, ql, qr, delta)
	}
	if qr > m {
		st.rangeAdd(v<<1|1, m, r, ql, qr, delta)
	}
	if st.min[v<<1] < st.min[v<<1|1] {
		st.min[v] = st.min[v<<1]
	} else {
		st.min[v] = st.min[v<<1|1]
	}
}

func (st *SegTree) AddRange(l, r, delta int) {
	if l < r {
		st.rangeAdd(1, 0, st.n, l, r, delta)
	}
}

func (st *SegTree) rangeMin(v, l, r, ql, qr int) int {
	if ql <= l && r <= qr {
		return st.min[v]
	}
	st.push(v)
	m := (l + r) >> 1
	res := int(1 << 60)
	if ql < m {
		val := st.rangeMin(v<<1, l, m, ql, qr)
		if val < res {
			res = val
		}
	}
	if qr > m {
		val := st.rangeMin(v<<1|1, m, r, ql, qr)
		if val < res {
			res = val
		}
	}
	return res
}

func (st *SegTree) QueryMin(l, r int) int {
	if l >= r {
		return int(1 << 60)
	}
	return st.rangeMin(1, 0, st.n, l, r)
}

func intersect(sx, sy, px, py, r float64) (float64, float64) {
	vx := px - sx
	vy := py - sy
	a := vx*vx + vy*vy
	b := 2 * (sx*vx + sy*vy)
	c := sx*sx + sy*sy - r*r
	disc := b*b - 4*a*c
	if disc < 0 {
		disc = 0
	}
	s := math.Sqrt(disc)
	t1 := (-b - s) / (2 * a)
	t2 := (-b + s) / (2 * a)
	t := t1
	if t < 0 || t > 1 {
		t = t2
	}
	if t < 0 || t > 1 {
		if t1 >= 0 && t1 <= 1 {
			t = t1
		} else if t2 >= 0 && t2 <= 1 {
			t = t2
		}
	}
	x := sx + t*vx
	y := sy + t*vy
	return x, y
}

func computeArc(sx, sy, r float64) (float64, float64) {
	xb, yb := intersect(sx, sy, r, 0, r)
	xa, ya := intersect(sx, sy, -r, 0, r)
	thB := math.Atan2(yb, xb)
	thA := math.Atan2(ya, xa)
	if thB < 0 {
		thB += 2 * math.Pi
	}
	if thA < 0 {
		thA += 2 * math.Pi
	}
	if thA < thB {
		thA += 2 * math.Pi
	}
	if thA > math.Pi {
		thA = math.Pi
	}
	if thB < 0 {
		thB = 0
	}
	return thB, thA
}

func lowerBound(a []float64, x float64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var r float64
	var n int
	if _, err := fmt.Fscan(reader, &r, &n); err != nil {
		return
	}

	events := make([]Event, n)
	angles := []float64{0, math.Pi}
	satID := 1
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(reader, &t)
		events[i].t = t
		if t == 1 {
			var x, y float64
			fmt.Fscan(reader, &x, &y)
			events[i].x = x
			events[i].y = y
			events[i].id = satID
			l, rr := computeArc(x, y, r)
			events[i].x = l
			events[i].y = rr
			angles = append(angles, l, rr)
			satID++
		} else if t == 2 {
			fmt.Fscan(reader, &events[i].i)
		} else if t == 3 {
			fmt.Fscan(reader, &events[i].i, &events[i].j)
		}
	}

	sort.Float64s(angles)
	uniq := angles[:0]
	for _, v := range angles {
		if len(uniq) == 0 || v-uniq[len(uniq)-1] > 1e-12 {
			uniq = append(uniq, v)
		}
	}
	angles = uniq
	seg := NewSegTree(len(angles))

	lIdx := make([]int, satID)
	rIdx := make([]int, satID)
	active := make([]bool, satID)

	for _, ev := range events {
		if ev.t == 1 {
			id := ev.id
			lIdx[id] = lowerBound(angles, ev.x)
			rIdx[id] = lowerBound(angles, ev.y)
		}
	}

	for _, ev := range events {
		switch ev.t {
		case 1:
			id := ev.id
			seg.AddRange(lIdx[id], rIdx[id], 1)
			active[id] = true
		case 2:
			id := ev.i
			if active[id] {
				seg.AddRange(lIdx[id], rIdx[id], -1)
				active[id] = false
			}
		case 3:
			id1 := ev.i
			id2 := ev.j
			l := math.Max(angles[lIdx[id1]], angles[lIdx[id2]])
			r := math.Min(angles[rIdx[id1]], angles[rIdx[id2]])
			if r <= l {
				fmt.Fprintln(writer, "NO")
				continue
			}
			L := lowerBound(angles, l)
			R := lowerBound(angles, r)
			if L >= R {
				fmt.Fprintln(writer, "NO")
				continue
			}
			m := seg.QueryMin(L, R)
			if m == 2 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
