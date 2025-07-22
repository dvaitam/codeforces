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

type edge struct {
	x, y int
}

type testCaseE struct {
	input    string
	expected []int64
}

const modE = 1000000007

func solveE(n int, edges []edge) []int64 {
	adj := make([][]int, n)
	for i := range adj {
		adj[i] = make([]int, n)
	}
	for _, e := range edges {
		adj[e.x-1][e.y-1] = 1
	}
	maxLen := 2 * n
	dp := make([][]int64, maxLen)
	for i := 0; i < maxLen; i++ {
		dp[i] = make([]int64, n)
	}
	for u := 0; u < n; u++ {
		dp[0][u] = 1
	}
	for s := 1; s < maxLen; s++ {
		for u := 0; u < n; u++ {
			var sum int64
			for v := 0; v < n; v++ {
				if adj[v][u] != 0 {
					sum += dp[s-1][v]
					if sum >= modE {
						sum -= modE
					}
				}
			}
			dp[s][u] = sum
		}
	}
	res := make([]int64, maxLen)
	for L := 1; L <= maxLen; L++ {
		if L == 1 {
			res[L-1] = 0
			continue
		}
		s := L - 1
		var total int64
		for u := 0; u < n; u++ {
			total += dp[s][u]
			if total >= modE {
				total -= modE
			}
		}
		res[L-1] = total
	}
	return res
}

func genCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(4) + 1
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges + 1)
	var edgesList []edge
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		for y == x {
			y = rng.Intn(n) + 1
		}
		edgesList = append(edgesList, edge{x, y})
		fmt.Fprintf(&sb, "%d %d 2 %d %d\n", x, y, x, y)
	}
	expected := solveE(n, edgesList)
	return testCaseE{input: sb.String(), expected: expected}
}

func runCaseE(bin string, tc testCaseE) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.expected) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.expected), len(fields))
	}
	for i, f := range fields {
		var v int64
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != tc.expected[i] {
			return fmt.Errorf("expected %v got %v", tc.expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
