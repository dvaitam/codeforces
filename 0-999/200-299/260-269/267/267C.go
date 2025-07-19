package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, l int
	if _, err := fmt.Fscan(in, &n, &l); err != nil {
		return
	}
	// adjacency
	v := make([][]int, n)
	d := make([][]float64, n)
	s := make([]int, l)
	e := make([]int, l)
	for k := 0; k < l; k++ {
		var i, j int
		var t float64
		fmt.Fscan(in, &i, &j, &t)
		i--
		j--
		s[k] = i
		e[k] = j
		v[i] = append(v[i], j)
		v[j] = append(v[j], i)
		d[i] = append(d[i], t)
		d[j] = append(d[j], t)
	}
	// build Laplacian-like counts with DFS
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, n)
	}
	u := make([]bool, n)
	var dfs func(int)
	dfs = func(i int) {
		u[i] = true
		for idx, j := range v[i] {
			a[i][j]++
			a[i][i]--
			if !u[j] {
				dfs(j)
			}
			// maintain same behavior as C++ (counts once per edge)
			_ = idx
		}
	}
	dfs(0)
	a[0][0]--
	// build matrix m: rows n-1, columns n (last is constant)
	rows := n - 1
	m := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		m[i] = make([]float64, n)
		for j := 0; j < rows; j++ {
			m[i][j] = a[i][j+1]
		}
		m[i][n-1] = -a[i][0]
	}
	y := make([]int, rows)
	eps := 1e-9
	singular := false
	// Gaussian elimination
	for i := 0; i < rows; i++ {
		if !u[i] {
			continue
		}
		// find pivot
		pivot := -1
		maxAbs := 0.0
		for t := 0; t < n; t++ {
			if abs := math.Abs(m[i][t]); abs > maxAbs {
				maxAbs = abs
				pivot = t
			}
		}
		if maxAbs < eps {
			singular = true
			break
		}
		y[i] = pivot
		// eliminate other rows
		for k := 0; k < rows; k++ {
			if k == i {
				continue
			}
			factor := m[k][pivot] / m[i][pivot]
			if factor == 0 {
				continue
			}
			for t := 0; t < n; t++ {
				m[k][t] -= m[i][t] * factor
			}
		}
	}
	if singular {
		fmt.Fprintln(out, "0.0")
		for i := 0; i < l; i++ {
			fmt.Fprintln(out, "0.0")
		}
		return
	}
	// recover x
	x := make([]float64, n)
	x[0] = 1.0
	for i := 0; i < rows; i++ {
		if !u[i] {
			continue
		}
		pivot := y[i]
		// x[index pivot+1] = m[i][n-1] / m[i][pivot]
		x[pivot+1] = m[i][n-1] / m[i][pivot]
	}
	// compute b
	b := 1e20
	for i := 0; i < n; i++ {
		if !u[i] {
			continue
		}
		for idx, j := range v[i] {
			z := math.Abs(x[j] - x[i])
			if z > eps {
				val := d[i][idx] / z
				if val < b {
					b = val
				}
			}
		}
	}
	// output b
	fmt.Fprintf(out, "%.13f\n", b)
	for k := 0; k < l; k++ {
		delta := (x[e[k]] - x[s[k]]) * b
		fmt.Fprintf(out, "%.13f\n", delta)
	}
}

func abs(v float64) float64 { return math.Abs(v) }
