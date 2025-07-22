package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"sort"
)

type interval struct{ l, r int64 }

func mergeIntervals(a []interval) []interval {
	if len(a) == 0 {
		return a
	}
	sort.Slice(a, func(i, j int) bool { return a[i].l < a[j].l })
	res := []interval{a[0]}
	for _, iv := range a[1:] {
		last := &res[len(res)-1]
		if iv.l <= last.r+1 {
			if iv.r > last.r {
				last.r = iv.r
			}
		} else {
			res = append(res, iv)
		}
	}
	return res
}

func contains(iv []interval, x int64) bool {
	// binary search
	i := sort.Search(len(iv), func(i int) bool { return iv[i].r >= x })
	if i < len(iv) && iv[i].l <= x && x <= iv[i].r {
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	const M int64 = 5000000000
	const MAX int64 = 2 * M
	x0, y0 := M, M
	var x, y int64 = x0, y0
	hor := make(map[int64][]interval)
	vert := make(map[int64][]interval)
	// record start point
	startX, startY := x0, y0
	for i := 0; i < n; i++ {
		var d byte
		var L int64
		fmt.Fscan(in, &d, &L)
		nx, ny := x, y
		switch d {
		case 'L':
			nx = x - L
			l, r := nx, x-1
			hor[y] = append(hor[y], interval{l, r})
		case 'R':
			nx = x + L
			l, r := x+1, nx
			hor[y] = append(hor[y], interval{l, r})
		case 'D':
			ny = y - L
			l, r := ny, y-1
			vert[x] = append(vert[x], interval{l, r})
		case 'U':
			ny = y + L
			l, r := y+1, ny
			vert[x] = append(vert[x], interval{l, r})
		}
		x, y = nx, ny
	}
	// merge intervals
	for ky, iv := range hor {
		hor[ky] = mergeIntervals(iv)
	}
	for kx, iv := range vert {
		vert[kx] = mergeIntervals(iv)
	}
	// prepare coords
	xs := make([]int64, 0, 4*(n+2))
	ys := make([]int64, 0, 4*(n+2))
	add := func(a *[]int64, v int64) {
		if v < 0 || v > MAX+1 {
			return
		}
		*a = append(*a, v)
	}
	add(&xs, 0)
	add(&xs, MAX+1)
	add(&ys, 0)
	add(&ys, MAX+1)
	// from horizontal segments
	for yy, ivs := range hor {
		add(&ys, yy)
		add(&ys, yy+1)
		for _, iv := range ivs {
			add(&xs, iv.l)
			add(&xs, iv.r+1)
		}
	}
	// from vertical segments
	for xx, ivs := range vert {
		add(&xs, xx)
		add(&xs, xx+1)
		for _, iv := range ivs {
			add(&ys, iv.l)
			add(&ys, iv.r+1)
		}
	}
	// start point
	add(&xs, startX)
	add(&xs, startX+1)
	add(&ys, startY)
	add(&ys, startY+1)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	ux := xs[:0]
	for i, v := range xs {
		if i == 0 || v != xs[i-1] {
			ux = append(ux, v)
		}
	}
	uy := ys[:0]
	for i, v := range ys {
		if i == 0 || v != ys[i-1] {
			uy = append(uy, v)
		}
	}
	xs, ys = ux, uy
	w, h := len(xs)-1, len(ys)-1
	// blocked and visited
	blocked := make([][]bool, w)
	vis := make([][]bool, w)
	for i := 0; i < w; i++ {
		blocked[i] = make([]bool, h)
		vis[i] = make([]bool, h)
	}
	// mark blocked cells
	for i := 0; i < w; i++ {
		dx := xs[i+1] - xs[i]
		for j := 0; j < h; j++ {
			dy := ys[j+1] - ys[j]
			if dx == 1 && dy == 1 {
				xi, yj := xs[i], ys[j]
				if xi == startX && yj == startY {
					blocked[i][j] = true
				} else if ivs, ok := hor[yj]; ok && contains(ivs, xi) {
					blocked[i][j] = true
				} else if ivs, ok := vert[xi]; ok && contains(ivs, yj) {
					blocked[i][j] = true
				}
			}
		}
	}
	// BFS from border
	type pair struct{ i, j int }
	q := make([]pair, 0, w*h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if i == 0 || j == 0 || i == w-1 || j == h-1 {
				if !blocked[i][j] && !vis[i][j] {
					vis[i][j] = true
					q = append(q, pair{i, j})
				}
			}
		}
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	invaded := big.NewInt(0)
	for bi := 0; bi < len(q); bi++ {
		u := q[bi]
		// add area
		dx := big.NewInt(xs[u.i+1] - xs[u.i])
		dy := big.NewInt(ys[u.j+1] - ys[u.j])
		area := big.NewInt(0).Mul(dx, dy)
		invaded.Add(invaded, area)
		for _, d := range dirs {
			ni, nj := u.i+d[0], u.j+d[1]
			if ni >= 0 && ni < w && nj >= 0 && nj < h && !blocked[ni][nj] && !vis[ni][nj] {
				vis[ni][nj] = true
				q = append(q, pair{ni, nj})
			}
		}
	}
	// total cells = (MAX+1)^2
	total := big.NewInt(0).Mul(big.NewInt(MAX+1), big.NewInt(MAX+1))
	safe := big.NewInt(0).Sub(total, invaded)
	fmt.Println(safe)
}
