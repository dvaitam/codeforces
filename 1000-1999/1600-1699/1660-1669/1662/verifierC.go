package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 998244353

type Matrix [][]int64

func makeMatrix(rows, cols int) Matrix {
	mat := make(Matrix, rows)
	for i := range mat {
		mat[i] = make([]int64, cols)
	}
	return mat
}

func matMul(a, b Matrix) Matrix {
	n := len(a)
	m := len(a[0])
	p := len(b[0])
	res := makeMatrix(n, p)
	for i := 0; i < n; i++ {
		for k := 0; k < m; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < p; j++ {
				res[i][j] = (res[i][j] + a[i][k]*b[k][j]) % mod
			}
		}
	}
	return res
}

func matPow(base Matrix, exp int) Matrix {
	n := len(base)
	res := makeMatrix(n, n)
	for i := 0; i < n; i++ {
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, base)
		}
		base = matMul(base, base)
		exp >>= 1
	}
	return res
}

func trace(mat Matrix) int64 {
	var s int64
	for i := 0; i < len(mat); i++ {
		s = (s + mat[i][i]) % mod
	}
	return s
}

func solveRef(input string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return 0, err
	}
	A := makeMatrix(n, n)
	deg := make([]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		A[u][v] = 1
		A[v][u] = 1
		deg[u]++
		deg[v]++
	}

	if k == 1 {
		return 0, nil
	}

	B1 := makeMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			B1[i][j] = A[i][j]
		}
	}
	B2 := matMul(A, A)
	for i := 0; i < n; i++ {
		B2[i][i] = (B2[i][i] - int64(deg[i]) + mod) % mod
	}
	if k == 2 {
		return trace(B2), nil
	}

	size := 2 * n
	M := makeMatrix(size, size)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			M[i][j] = A[i][j]
		}
		M[i][n+i] = (mod - int64(deg[i]-1)%mod) % mod
		M[n+i][i] = 1
	}

	P := matPow(M, k-2)

	X := makeMatrix(size, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			X[i][j] = B2[i][j]
			X[n+i][j] = B1[i][j]
		}
	}

	Y := matMul(P, X)
	Bk := Y[:n]
	return trace(Bk), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (int64, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int64
	if _, err := fmt.Sscan(out, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val % mod, nil
}

type testCase struct {
	name  string
	input string
}

func makeCase(name string, n, m, k int, edges [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{name: name, input: sb.String()}
}

func randomGraph(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rng.Intn(2) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	if len(edges) == 0 {
		edges = append(edges, [2]int{1, n})
	}
	return edges
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("path3_k2", 3, 2, 2, [][2]int{{1, 2}, {2, 3}}),
		makeCase("triangle_k3", 3, 3, 3, [][2]int{{1, 2}, {2, 3}, {1, 3}}),
		makeCase("square", 4, 4, 4, [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 30; i++ {
		n := rng.Intn(5) + 3
		k := rng.Intn(5) + 1
		edges := randomGraph(n, rng)
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), n, len(edges), k, edges))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%d\nactual:%d\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
