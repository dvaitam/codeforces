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

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solve(n, m, k int, g [][]int64) string {
	size := n + m + 5
	a := make([][]int64, size)
	for i := range a {
		a[i] = make([]int64, size)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			v := g[i-1][j-1]
			x := i + j
			y := i - j + m
			a[x][y] = v
		}
	}
	limit := n + m
	for i := 1; i <= limit; i++ {
		row := a[i]
		prev := a[i-1]
		for j := 1; j <= limit; j++ {
			row[j] += prev[j] + row[j-1] - prev[j-1]
		}
	}
	sum := func(l, L, r, R int) int64 {
		return a[r][R] - a[l-1][R] - a[r][L-1] + a[l-1][L-1]
	}
	var best int64 = -1
	ansX, ansY := 0, 0
	for i := k; i <= n-k+1; i++ {
		for j := k; j <= m-k+1; j++ {
			x := i + j
			y := i - j + m
			var s int64
			for z := 0; z < k; z++ {
				s += sum(x-z, y-z, x+z, y+z)
			}
			if s > best {
				best = s
				ansX, ansY = i, j
			}
		}
	}
	return fmt.Sprintf("%d %d", ansX, ansY)
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		kMax := (min(n, m) + 1) / 2
		if kMax == 0 {
			kMax = 1
		}
		k := rand.Intn(kMax) + 1
		g := make([][]int64, n)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			g[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				val := int64(rand.Intn(11))
				g[i][j] = val
				if j > 0 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", val)
			}
			buf.WriteByte('\n')
		}
		expected := solve(n, m, k, g)
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
