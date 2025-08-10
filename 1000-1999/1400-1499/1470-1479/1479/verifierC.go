package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// runCmd executes the binary at cmdPath with the provided input and returns
// its combined standard output and error.
func runCmd(cmdPath, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, cmdPath)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	return out.String(), err
}

// Edge represents a directed weighted edge.
type Edge struct {
	to, w int
}

// verify parses the candidate output and checks that the graph encodes all
// integers in the range [l, r]. It also ensures there is no path whose length
// is less than l.
func verify(l, r int, out string) error {
	rd := strings.NewReader(out)
	var verdict string
	if _, err := fmt.Fscan(rd, &verdict); err != nil {
		return fmt.Errorf("failed to read verdict: %v", err)
	}
	if verdict != "YES" {
		return fmt.Errorf("expected YES, got %q", verdict)
	}

	var n, m int
	if _, err := fmt.Fscan(rd, &n, &m); err != nil {
		return fmt.Errorf("failed to read n and m: %v", err)
	}
	if n <= 0 || m < 0 {
		return fmt.Errorf("invalid n or m")
	}

	edges := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var a, b, c int
		if _, err := fmt.Fscan(rd, &a, &b, &c); err != nil {
			return fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
		if a < 1 || b < 1 || a >= b || a > n || b > n || c < 0 {
			return fmt.Errorf("invalid edge %d", i+1)
		}
		edges[a] = append(edges[a], Edge{to: b, w: c})
	}

	// Dynamic programming over the DAG to compute reachable path sums.
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, r+1)
	}
	dp[1][0] = true
	for u := 1; u <= n; u++ {
		for sum := 0; sum <= r; sum++ {
			if !dp[u][sum] {
				continue
			}
			for _, e := range edges[u] {
				if sum+e.w <= r {
					dp[e.to][sum+e.w] = true
				}
			}
		}
	}

	for x := 0; x < l; x++ {
		if dp[n][x] {
			return fmt.Errorf("found path with sum %d < %d", x, l)
		}
	}
	for x := l; x <= r; x++ {
		if !dp[n][x] {
			return fmt.Errorf("missing path with sum %d", x)
		}
	}
	return nil
}

func genTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 100)
	for i := range tests {
		l := rng.Intn(20) + 1
		r := l + rng.Intn(20)
		tests[i] = fmt.Sprintf("%d %d\n", l, r)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := genTests()
	for i, t := range tests {
		var l, r int
		fmt.Sscan(t, &l, &r)
		got, err := runCmd(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(l, r, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\n", i+1, err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
