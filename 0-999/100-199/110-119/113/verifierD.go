package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(n, m, A, B int, edges [][2]int, p []float64) []float64 {
	A--
	B--
	if A < B {
		A, B = B, A
	}
	g := make([][]int, n)
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	id := make([][]int, n)
	for i := 0; i < n; i++ {
		id[i] = make([]int, n)
	}
	N := 0
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			id[i][j] = N
			N++
		}
	}
	a := make([][]float64, N)
	for i := 0; i < N; i++ {
		a[i] = make([]float64, N)
	}
	edge := func(ea, b1, c, d int, e float64) {
		if ea < b1 {
			ea, b1 = b1, ea
		}
		if c < d {
			c, d = d, c
		}
		i1 := id[ea][b1]
		i2 := id[c][d]
		a[i1][i2] += e
	}
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			if i == j {
				edge(i, j, i, j, 1)
				continue
			}
			edge(i, j, i, j, p[i]*p[j])
			giLen := float64(len(g[i]))
			gjLen := float64(len(g[j]))
			for _, qi := range g[i] {
				edge(i, j, qi, j, (1-p[i])/giLen*p[j])
				for _, qj := range g[j] {
					edge(i, j, qi, qj, (1-p[i])/giLen*(1-p[j])/gjLen)
				}
			}
			for _, qj := range g[j] {
				edge(i, j, i, qj, (1-p[j])/gjLen*p[i])
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			idx := id[i][j]
			a[idx][idx] -= 1
		}
	}
	l := make([][]float64, N)
	u := make([][]float64, N)
	for i := 0; i < N; i++ {
		l[i] = make([]float64, N)
		u[i] = make([]float64, N)
	}
	for j := 0; j < N; j++ {
		u[0][j] = a[0][j]
	}
	for j := 1; j < N; j++ {
		l[j][0] = a[j][0] / u[0][0]
	}
	for i := 1; i < N; i++ {
		for j := i; j < N; j++ {
			sum := a[i][j]
			for k := 0; k < i; k++ {
				sum -= l[i][k] * u[k][j]
			}
			u[i][j] = sum
		}
		for j := i + 1; j < N; j++ {
			sum := a[j][i]
			for k := 0; k < i; k++ {
				sum -= l[j][k] * u[k][i]
			}
			l[j][i] = sum / u[i][i]
		}
	}
	for i := 0; i < N; i++ {
		l[i][i] = 1
	}
	x := make([]float64, N)
	y := make([]float64, N)
	bvec := make([]float64, N)
	res := make([]float64, n)
	idxAB := id[A][B]
	for d := 0; d < n; d++ {
		for i := 0; i < N; i++ {
			bvec[i] = 0
		}
		bvec[id[d][d]] = 1
		for i := 0; i < N; i++ {
			sum := bvec[i]
			for j := 0; j < i; j++ {
				sum -= l[i][j] * y[j]
			}
			y[i] = sum
		}
		for i := N - 1; i >= 0; i-- {
			sum := y[i]
			for j := i + 1; j < N; j++ {
				sum -= u[i][j] * x[j]
			}
			x[i] = sum / u[i][i]
		}
		res[d] = x[idxAB]
	}
	return res
}

func formatRes(res []float64) string {
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.8f", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	edges := make([][2]int, m)
	used := make(map[[2]int]bool)
	idx := 0
	for idx < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u < v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		edges[idx] = key
		idx++
	}
	A := rng.Intn(n) + 1
	B := rng.Intn(n) + 1
	for B == A {
		B = rng.Intn(n) + 1
	}
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		probs[i] = rng.Float64()*0.8 + 0.1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, A, B))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%.2f\n", probs[i]))
	}
	return sb.String(), formatRes(solveD(n, m, A, B, edges, probs))
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
