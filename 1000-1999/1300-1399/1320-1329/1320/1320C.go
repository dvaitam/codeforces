package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	maxValue = 1000005
	inf      = int64(1 << 60)
)

type weapon struct {
	attack int
	cost   int64
}

type monster struct {
	x, y int
	z    int64
}

type segTree struct {
	n    int
	max  []int64
	lazy []int64
}

func newSegTree(values []int64) *segTree {
	st := &segTree{
		n:    len(values),
		max:  make([]int64, len(values)*4),
		lazy: make([]int64, len(values)*4),
	}
	st.build(1, 0, st.n-1, values)
	return st
}

func (st *segTree) build(node, l, r int, values []int64) {
	if l == r {
		st.max[node] = values[l]
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, values)
	st.build(node<<1|1, mid+1, r, values)
	st.pull(node)
}

func (st *segTree) pull(node int) {
	if st.max[node<<1] > st.max[node<<1|1] {
		st.max[node] = st.max[node<<1]
	} else {
		st.max[node] = st.max[node<<1|1]
	}
}

func (st *segTree) push(node int) {
	if st.lazy[node] == 0 {
		return
	}
	val := st.lazy[node]
	st.max[node<<1] += val
	st.lazy[node<<1] += val
	st.max[node<<1|1] += val
	st.lazy[node<<1|1] += val
	st.lazy[node] = 0
}

func (st *segTree) updateRange(l, r int, val int64) {
	if l < 0 {
		l = 0
	}
	if r >= st.n {
		r = st.n - 1
	}
	if l > r {
		return
	}
	st.update(1, 0, st.n-1, l, r, val)
}

func (st *segTree) update(node, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.max[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.update(node<<1, l, mid, ql, qr, val)
	st.update(node<<1|1, mid+1, r, ql, qr, val)
	st.pull(node)
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	n := fs.nextInt()
	m := fs.nextInt()
	p := fs.nextInt()

	bestWeapon := make([]int64, maxValue)
	for i := range bestWeapon {
		bestWeapon[i] = inf
	}

	for i := 0; i < n; i++ {
		a := fs.nextInt()
		c := fs.nextInt()
		cost := int64(c)
		if cost < bestWeapon[a] {
			bestWeapon[a] = cost
		}
	}

	base := make([]int64, maxValue)
	for i := range base {
		base[i] = -inf
	}

	for i := 0; i < m; i++ {
		b := fs.nextInt()
		c := fs.nextInt()
		val := -int64(c)
		if val > base[b] {
			base[b] = val
		}
	}

	monsters := make([]monster, p)
	for i := 0; i < p; i++ {
		x := fs.nextInt()
		y := fs.nextInt()
		z := fs.nextInt()
		monsters[i] = monster{x: x, y: y, z: int64(z)}
	}
	sort.Slice(monsters, func(i, j int) bool {
		if monsters[i].x == monsters[j].x {
			return monsters[i].y < monsters[j].y
		}
		return monsters[i].x < monsters[j].x
	})

	st := newSegTree(base)

	var weapons []weapon
	for attack := 1; attack < maxValue; attack++ {
		if bestWeapon[attack] < inf {
			weapons = append(weapons, weapon{attack: attack, cost: bestWeapon[attack]})
		}
	}

	ans := int64(-inf)
	monsterIdx := 0

	for _, w := range weapons {
		for monsterIdx < len(monsters) && monsters[monsterIdx].x < w.attack {
			y := monsters[monsterIdx].y
			st.updateRange(y+1, maxValue-1, monsters[monsterIdx].z)
			monsterIdx++
		}
		current := st.max[1] - w.cost
		if current > ans {
			ans = current
		}
	}

	fmt.Fprintln(writer, ans)
}
