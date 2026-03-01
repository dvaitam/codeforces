package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	k int
}

func possible(n, k int) bool {
	if k < 2 || k > 3 || n <= 3 {
		return false
	}
	return !(n == 4 && k == 2)
}

func diameter(adj [][]bool) int {
	n := len(adj)
	diam := 0
	q := make([]int, n)
	dist := make([]int, n)
	for s := 0; s < n; s++ {
		for i := 0; i < n; i++ {
			dist[i] = -1
		}
		head, tail := 0, 0
		q[tail] = s
		tail++
		dist[s] = 0
		for head < tail {
			v := q[head]
			head++
			for to := 0; to < n; to++ {
				if !adj[v][to] || dist[to] != -1 {
					continue
				}
				dist[to] = dist[v] + 1
				q[tail] = to
				tail++
			}
		}
		for i := 0; i < n; i++ {
			if dist[i] == -1 {
				return -1
			}
			if dist[i] > diam {
				diam = dist[i]
			}
		}
	}
	return diam
}

func validateOutput(n, k int, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	if tokens[0] == "-1" {
		if len(tokens) != 1 {
			return fmt.Errorf("-1 must be the only token")
		}
		if possible(n, k) {
			return fmt.Errorf("reported impossible for a possible case")
		}
		return nil
	}
	if !possible(n, k) {
		return fmt.Errorf("reported a graph for an impossible case")
	}

	m, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid edge count: %v", err)
	}
	if m < 0 || m > n*(n-1)/2 {
		return fmt.Errorf("invalid edge count %d", m)
	}
	if len(tokens) != 1+2*m {
		return fmt.Errorf("expected %d edge endpoint tokens, got %d", 2*m, len(tokens)-1)
	}

	red := make([][]bool, n)
	for i := 0; i < n; i++ {
		red[i] = make([]bool, n)
	}

	for i := 0; i < m; i++ {
		a, errA := strconv.Atoi(tokens[1+2*i])
		b, errB := strconv.Atoi(tokens[2+2*i])
		if errA != nil || errB != nil {
			return fmt.Errorf("invalid edge on line %d", i+1)
		}
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("vertex out of range in edge (%d, %d)", a, b)
		}
		if a == b {
			return fmt.Errorf("self-loop (%d, %d)", a, b)
		}
		a--
		b--
		if red[a][b] {
			return fmt.Errorf("duplicate edge (%d, %d)", a+1, b+1)
		}
		red[a][b], red[b][a] = true, true
	}

	white := make([][]bool, n)
	for i := 0; i < n; i++ {
		white[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			white[i][j] = !red[i][j]
			white[j][i] = white[i][j]
		}
	}

	dr := diameter(red)
	dw := diameter(white)
	if min(dr, dw) != k {
		return fmt.Errorf("wrong colorfulness: dr=%d, dw=%d, min=%d, want=%d", dr, dw, min(dr, dw), k)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateTestCases() []testCase {
	tests := make([]testCase, 0)
	for n := 2; n <= 30; n++ {
		for k := 1; k <= 8; k++ {
			tests = append(tests, testCase{n: n, k: k})
		}
	}
	for _, n := range []int{50, 100, 250, 500, 1000} {
		tests = append(tests,
			testCase{n: n, k: 1},
			testCase{n: n, k: 2},
			testCase{n: n, k: 3},
			testCase{n: n, k: 4},
			testCase{n: n, k: 1000},
		)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	tests := generateTestCases()

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d (n=%d, k=%d): runtime error: %v\nstderr: %s\n", i+1, tc.n, tc.k, err, stderr.String())
			os.Exit(1)
		}
		if err := validateOutput(tc.n, tc.k, out.String()); err != nil {
			fmt.Printf("test %d (n=%d, k=%d) failed: %v\nprogram output:\n%s\n", i+1, tc.n, tc.k, err, out.String())
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
