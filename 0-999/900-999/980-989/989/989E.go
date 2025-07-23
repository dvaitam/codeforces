package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

type Line struct {
	points []int
}

func matMul(a, b [][]float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := 0; i < n; i++ {
		c[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func matVecMul(a [][]float64, v []float64) []float64 {
	n := len(a)
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := 0.0
		row := a[i]
		for j := 0; j < n; j++ {
			sum += row[j] * v[j]
		}
		res[i] = sum
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &y[i])
	}

	type key struct{ dx, dy, c int }
	lineMap := make(map[key]map[int]struct{})

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := x[j] - x[i]
			dy := y[j] - y[i]
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			c := dy*x[i] - dx*y[i]
			k := key{dx, dy, c}
			if lineMap[k] == nil {
				lineMap[k] = make(map[int]struct{})
			}
			lineMap[k][i] = struct{}{}
			lineMap[k][j] = struct{}{}
		}
	}

	lines := make([]Line, 0, len(lineMap))
	for _, m := range lineMap {
		pts := make([]int, 0, len(m))
		for p := range m {
			pts = append(pts, p)
		}
		if len(pts) >= 2 {
			lines = append(lines, Line{pts})
		}
	}

	pointLines := make([][]int, n)
	for idx, ln := range lines {
		for _, p := range ln.points {
			pointLines[p] = append(pointLines[p], idx)
		}
	}

	P := make([][]float64, n)
	for i := 0; i < n; i++ {
		P[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		deg := len(pointLines[i])
		if deg == 0 {
			continue
		}
		for _, li := range pointLines[i] {
			pts := lines[li].points
			prob := 1.0 / float64(deg) / float64(len(pts))
			for _, j := range pts {
				P[i][j] += prob
			}
		}
	}
	maxM := 10000

	mats := make([][][]float64, 0)
	mats = append(mats, P)
	for (1 << uint(len(mats))) <= maxM {
		next := matMul(mats[len(mats)-1], mats[len(mats)-1])
		mats = append(mats, next)
	}

	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var t, m int
		fmt.Fscan(in, &t, &m)
		t--
		step := m - 1
		w := make([]float64, n)
		w[t] = 1
		bit := 0
		for step > 0 {
			if step&1 == 1 {
				w = matVecMul(mats[bit], w)
			}
			step >>= 1
			bit++
		}
		ans := 0.0
		for _, ln := range lines {
			sum := 0.0
			for _, p := range ln.points {
				sum += w[p]
			}
			prob := sum / float64(len(ln.points))
			if prob > ans {
				ans = prob
			}
		}
		fmt.Fprintf(out, "%.9f\n", ans)
	}
}
