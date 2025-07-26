package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Line for convex hull
type Line struct{ m, b, xLeft int64 }

// Hull represents convex hull for max queries
type Hull struct{ lines []Line }

// intersect returns ceil((b1-b2)/(m2-m1))
func intersect(l1, l2 Line) int64 {
	// m1 < m2
	num := l1.b - l2.b
	den := l2.m - l1.m
	if den < 0 {
		num = -num
		den = -den
	}
	if num <= 0 {
		return 0
	}
	return (num + den - 1) / den
}

// build builds hull from sorted lines by slope
func (h *Hull) build(lines []Line) {
	h.lines = h.lines[:0]
	for _, l := range lines {
		if len(h.lines) > 0 && h.lines[len(h.lines)-1].m == l.m {
			// keep higher b
			if h.lines[len(h.lines)-1].b >= l.b {
				continue
			}
			h.lines = h.lines[:len(h.lines)-1]
		}
		for len(h.lines) >= 2 {
			x := intersect(h.lines[len(h.lines)-2], h.lines[len(h.lines)-1])
			x2 := intersect(h.lines[len(h.lines)-2], l)
			if x2 <= x {
				h.lines = h.lines[:len(h.lines)-1]
			} else {
				break
			}
		}
		h.lines = append(h.lines, l)
	}
	// set xLeft
	for i := range h.lines {
		if i == 0 {
			h.lines[i].xLeft = 0
		} else {
			h.lines[i].xLeft = intersect(h.lines[i-1], h.lines[i])
		}
	}
}

// query returns max m*x + b at x
func (h *Hull) query(x int64) int64 {
	if len(h.lines) == 0 {
		return 0
	}
	// binary search for last line with xLeft <= x
	lo, hi := 0, len(h.lines)-1
	for lo < hi {
		mid := (lo + hi + 1) >> 1
		if h.lines[mid].xLeft <= x {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	l := h.lines[lo]
	return l.m*x + l.b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, q int
	fmt.Fscan(reader, &n, &q)
	p := make([]int, n+1)
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
		children[p[i]] = append(children[p[i]], i)
	}
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	flatA := make([]int64, n)
	flatC := make([]int64, n)
	timer := 0
	type Item struct {
		u, idx     int
		sumA, sumB int64
	}
	stack := []Item{{1, 0, a[1], b[1]}}
	for len(stack) > 0 {
		itm := &stack[len(stack)-1]
		u := itm.u
		if itm.idx == 0 {
			tin[u] = timer
			flatA[timer] = itm.sumA
			if itm.sumB < 0 {
				flatC[timer] = -itm.sumB
			} else {
				flatC[timer] = itm.sumB
			}
			timer++
		}
		if itm.idx < len(children[u]) {
			v := children[u][itm.idx]
			itm.idx++
			stack = append(stack, Item{v, 0, itm.sumA + a[v], itm.sumB + b[v]})
		} else {
			tout[u] = timer - 1
			stack = stack[:len(stack)-1]
		}
	}
	// sqrt decomposition
	const Bsize = 450
	nb := (n + Bsize - 1) / Bsize
	type Block struct {
		l, r         int
		add          int64
		A, C         []int64
		hullG, hullH Hull
		built        bool
	}
	blocks := make([]*Block, nb)
	for i := 0; i < nb; i++ {
		l := i * Bsize
		r := (i + 1) * Bsize
		if r > n {
			r = n
		}
		blk := &Block{l: l, r: r, add: 0, built: false}
		blk.A = make([]int64, r-l)
		blk.C = make([]int64, r-l)
		for j := l; j < r; j++ {
			blk.A[j-l] = flatA[j]
			blk.C[j-l] = flatC[j]
		}
		blocks[i] = blk
	}
	// helpers
	maxI := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	minI := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	abs64 := func(x int64) int64 {
		if x < 0 {
			return -x
		}
		return x
	}
	// build block hulls
	buildBlock := func(blk *Block) {
		var linesG, linesH []Line
		for i := range blk.A {
			D := blk.C[i]
			A0 := blk.A[i]
			linesG = append(linesG, Line{m: D, b: D * A0})
			linesH = append(linesH, Line{m: -D, b: -D * A0})
		}
		sort.Slice(linesG, func(i, j int) bool { return linesG[i].m < linesG[j].m })
		sort.Slice(linesH, func(i, j int) bool { return linesH[i].m < linesH[j].m })
		blk.hullG.build(linesG)
		blk.hullH.build(linesH)
		blk.built = true
	}
	// process queries
	for qi := 0; qi < q; qi++ {
		var typ, v int
		fmt.Fscan(reader, &typ, &v)
		l := tin[v]
		r := tout[v]
		if typ == 1 {
			var x int64
			fmt.Fscan(reader, &x)
			for _, blk := range blocks {
				if r < blk.l || blk.r-1 < l {
					continue
				}
				if l <= blk.l && blk.r-1 <= r {
					blk.add += x
				} else {
					for j := maxI(l, blk.l); j <= minI(r, blk.r-1); j++ {
						blk.A[j-blk.l] += x
					}
					blk.built = false
				}
			}
		} else {
			ans := int64(0)
			for _, blk := range blocks {
				if r < blk.l || blk.r-1 < l {
					continue
				}
				if l <= blk.l && blk.r-1 <= r {
					if !blk.built {
						buildBlock(blk)
					}
					t := blk.add
					v1 := blk.hullG.query(t)
					v2 := blk.hullH.query(t)
					if v1 > ans {
						ans = v1
					}
					if v2 > ans {
						ans = v2
					}
				} else {
					for j := maxI(l, blk.l); j <= minI(r, blk.r-1); j++ {
						vj := blk.C[j-blk.l] * abs64(blk.A[j-blk.l]+blk.add)
						if vj > ans {
							ans = vj
						}
					}
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
