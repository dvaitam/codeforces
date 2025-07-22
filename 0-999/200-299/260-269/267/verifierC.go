package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct {
	u, v int
	t    float64
}

type testCase struct {
	input string
	ans   []float64
}

func solve(n int, edges []edge) []float64 {
	l := len(edges)
	v := make([][]int, n)
	d := make([][]float64, n)
	s := make([]int, l)
	e := make([]int, l)
	for k, ed := range edges {
		i := ed.u
		j := ed.v
		t := ed.t
		s[k] = i
		e[k] = j
		v[i] = append(v[i], j)
		v[j] = append(v[j], i)
		d[i] = append(d[i], t)
		d[j] = append(d[j], t)
	}
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
			_ = idx
		}
	}
	dfs(0)
	a[0][0]--
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
	for i := 0; i < rows; i++ {
		if !u[i] {
			continue
		}
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
	res := make([]float64, 1+l)
	if singular {
		return res
	}
	x := make([]float64, n)
	x[0] = 1.0
	for i := 0; i < rows; i++ {
		if !u[i] {
			continue
		}
		pivot := y[i]
		x[pivot+1] = m[i][n-1] / m[i][pivot]
	}
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
	res[0] = b
	for k := 0; k < l; k++ {
		delta := (x[e[k]] - x[s[k]]) * b
		res[k+1] = delta
	}
	return res
}

func buildCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 2
	l := rng.Intn(5) + 1
	edges := make([]edge, l)
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, l)
	for i := 0; i < l; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for v == u {
			v = rng.Intn(n)
		}
		t := rng.Float64()*9 + 1
		fmt.Fprintf(&in, "%d %d %.6f\n", u+1, v+1, t)
		edges[i] = edge{u, v, t}
	}
	ans := solve(n, edges)
	return testCase{in.String(), ans}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.ans) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.ans), len(fields))
	}
	for i, f := range fields {
		var val float64
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad float: %v", err)
		}
		if math.Abs(val-tc.ans[i]) > 1e-6 {
			return fmt.Errorf("mismatch at %d: expected %.6f got %.6f", i, tc.ans[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		cases[i] = buildCase(rng)
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
