package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveD(n, m, a, b int, edges [][2]int, p []float64) []float64 {
	a--
	b--

	adj := make([][]int, n)
	deg := make([]int, n)
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}

	ans := make([]float64, n)
	if a == b {
		ans[a] = 1.0
		return ans
	}

	P := make([][]float64, n)
	for i := 0; i < n; i++ {
		P[i] = make([]float64, n)
		P[i][i] = p[i]
		for _, v := range adj[i] {
			P[i][v] = (1.0 - p[i]) / float64(deg[i])
		}
	}

	stateIdx := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		stateIdx[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i != j {
				stateIdx[i][j] = idx
				idx++
			}
		}
	}

	N := idx
	A := make([][]float64, N)
	for i := 0; i < N; i++ {
		A[i] = make([]float64, N)
	}
	B := make([]float64, N)

	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			if x == y {
				continue
			}
			Sp := stateIdx[x][y]
			for u := 0; u < n; u++ {
				if P[x][u] == 0 {
					continue
				}
				for v := 0; v < n; v++ {
					if P[y][v] == 0 {
						continue
					}
					if u != v {
						S := stateIdx[u][v]
						A[S][Sp] -= P[x][u] * P[y][v]
					}
				}
			}
		}
	}

	for i := 0; i < N; i++ {
		A[i][i] += 1.0
	}
	B[stateIdx[a][b]] = 1.0

	// Gaussian elimination with partial pivoting
	for i := 0; i < N; i++ {
		pivot := i
		for j := i + 1; j < N; j++ {
			if math.Abs(A[j][i]) > math.Abs(A[pivot][i]) {
				pivot = j
			}
		}

		A[i], A[pivot] = A[pivot], A[i]
		B[i], B[pivot] = B[pivot], B[i]

		inv := 1.0 / A[i][i]
		for j := i; j < N; j++ {
			A[i][j] *= inv
		}
		B[i] *= inv

		for j := 0; j < N; j++ {
			if i != j && A[j][i] != 0 {
				factor := A[j][i]
				for k := i; k < N; k++ {
					A[j][k] -= factor * A[i][k]
				}
				B[j] -= factor * B[i]
			}
		}
	}

	for k := 0; k < n; k++ {
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				if x != y {
					S := stateIdx[x][y]
					ans[k] += B[S] * P[x][k] * P[y][k]
				}
			}
		}
	}

	return ans
}

func formatRes(res []float64) string {
	parts := make([]string, len(res))
	for i, v := range res {
		parts[i] = fmt.Sprintf("%.9f", v)
	}
	return strings.Join(parts, " ")
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
	input := sb.String()
	// Round probs to 2 decimals to match what the binary reads from stdin
	roundedProbs := make([]float64, n)
	for i := 0; i < n; i++ {
		roundedProbs[i], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", probs[i]), 64)
	}
	res := solveD(n, m, A, B, edges, roundedProbs)
	return input, formatRes(res)
}

func parseFloats(s string) []float64 {
	parts := strings.Fields(s)
	res := make([]float64, len(parts))
	for i, p := range parts {
		res[i], _ = strconv.ParseFloat(p, 64)
	}
	return res
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

	gotVals := parseFloats(got)
	expVals := parseFloats(exp)
	if len(gotVals) != len(expVals) {
		return fmt.Errorf("expected %d values got %d\nexpected %q\ngot %q", len(expVals), len(gotVals), exp, got)
	}
	for i := range gotVals {
		if math.Abs(gotVals[i]-expVals[i]) > 1e-5 {
			return fmt.Errorf("value %d: expected %.9f got %.9f\nexpected %q\ngot %q", i, expVals[i], gotVals[i], exp, got)
		}
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
