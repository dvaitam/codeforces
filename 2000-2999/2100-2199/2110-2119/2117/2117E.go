package main

import (
	"bufio"
	"fmt"
	"os"
)

func solveCase(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	x := make([]int, n+1)
	y := make([]int, n+1)
	occX := make([][]int, n+1)
	occY := make([][]int, n+1)

	for i := 1; i <= n; i++ {
		if i&1 == 1 {
			x[i] = a[i]
			y[i] = b[i]
		} else {
			x[i] = b[i]
			y[i] = a[i]
		}
		occX[x[i]] = append(occX[x[i]], i)
		occY[y[i]] = append(occY[y[i]], i)
	}

	ans := 0
	for v := 1; v <= n; v++ {
		lx, ly := 0, 0
		if len(occX[v]) > 0 {
			lx = occX[v][len(occX[v])-1]
		}
		if len(occY[v]) > 0 {
			ly = occY[v][len(occY[v])-1]
		}
		if lx < ly {
			if lx > ans {
				ans = lx
			}
		} else {
			if ly > ans {
				ans = ly
			}
		}
	}

	px := make([]int, n+1)
	py := make([]int, n+1)
	ptrX := make([]int, n+1)
	ptrY := make([]int, n+1)

	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]int, size<<1)

	calc := func(v int) int {
		sx, sy := 0, 0
		if ptrX[v] >= 0 {
			sx = occX[v][ptrX[v]]
		}
		if ptrY[v] >= 0 {
			sy = occY[v][ptrY[v]]
		}
		lx := px[v]
		if sy > 0 {
			t := sy - 1
			if t > lx {
				lx = t
			}
		}
		ly := py[v]
		if sx > 0 {
			t := sx - 1
			if t > ly {
				ly = t
			}
		}
		if lx < ly {
			return lx
		}
		return ly
	}

	update := func(v int) {
		p := size + v - 1
		seg[p] = calc(v)
		for p >>= 1; p > 0; p >>= 1 {
			if seg[p<<1] >= seg[p<<1|1] {
				seg[p] = seg[p<<1]
			} else {
				seg[p] = seg[p<<1|1]
			}
		}
	}

	for v := 1; v <= n; v++ {
		ptrX[v] = len(occX[v]) - 1
		for ptrX[v] >= 0 && occX[v][ptrX[v]] <= 1 {
			ptrX[v]--
		}
		ptrY[v] = len(occY[v]) - 1
		for ptrY[v] >= 0 && occY[v][ptrY[v]] <= 1 {
			ptrY[v]--
		}
		seg[size+v-1] = calc(v)
	}
	for i := size - 1; i >= 1; i-- {
		if seg[i<<1] >= seg[i<<1|1] {
			seg[i] = seg[i<<1]
		} else {
			seg[i] = seg[i<<1|1]
		}
	}

	if seg[1] > ans {
		ans = seg[1]
	}

	for r := 1; r < n; r++ {
		v := x[r]
		px[v] = r
		update(v)

		v = y[r]
		py[v] = r
		update(v)

		v = x[r+1]
		for ptrX[v] >= 0 && occX[v][ptrX[v]] <= r+1 {
			ptrX[v]--
		}
		update(v)

		v = y[r+1]
		for ptrY[v] >= 0 && occY[v][ptrY[v]] <= r+1 {
			ptrY[v]--
		}
		update(v)

		if seg[1] > ans {
			ans = seg[1]
		}
	}

	fmt.Fprintln(out, ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		solveCase(in, out)
	}
}
