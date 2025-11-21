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

type testCase struct {
	name  string
	input string
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func solveRef(input string) (float64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return 0, err
	}
	c := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
		total += c[i]
	}

	C := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]float64, i+1)
		C[i][0], C[i][i] = 1.0, 1.0
		for j := 1; j < i; j++ {
			C[i][j] = C[i-1][j] + C[i-1][j-1]
		}
	}

	dp := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]float64, total+1)
	}
	dp[0][0] = 1.0
	now := 0
	for i := 0; i < n; i++ {
		ci := c[i]
		now += ci
		for j := i + 1; j >= 1; j-- {
			for p := now; p >= ci; p-- {
				dp[j][p] += dp[j-1][p-ci]
			}
		}
	}

	var ans float64
	nf := float64(n)
	xf := float64(x)
	for i := 1; i <= n; i++ {
		limit := (nf/float64(i) + 1.0) * xf / 2.0
		invCi := 1.0 / C[n][i]
		for j := 1; j <= total; j++ {
			ways := dp[i][j]
			if ways == 0 {
				continue
			}
			value := float64(j) / float64(i)
			cost := minFloat(limit, value)
			ans += ways * invCi * cost
		}
	}
	return ans, nil
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

func parseOutput(out string) (float64, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val float64
	if _, err := fmt.Sscan(out, &val); err != nil {
		return 0, fmt.Errorf("failed to parse float: %v", err)
	}
	return val, nil
}

func makeCase(name string, n, x int, costs []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", costs[i])
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 80; i++ {
		n := rng.Intn(6) + 1
		total := 0
		costs := make([]int, n)
		for j := 0; j < n; j++ {
			costs[j] = rng.Intn(20) + 1
			total += costs[j]
		}
		x := rng.Intn(20) + 1
		for _, c := range costs {
			if c < x {
				// ensure constraint x <= ci
				x = c
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), n, x, costs))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single", 1, 5, []int{5}),
		makeCase("two_simple", 2, 20, []int{25, 30}),
		makeCase("three", 3, 10, []int{10, 15, 20}),
		makeCase("equal_costs", 4, 12, []int{12, 12, 12, 12}),
		makeCase("increasing_costs", 5, 8, []int{8, 9, 10, 11, 12}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	const eps = 1e-6
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
		diff := got - expect
		if diff < 0 {
			diff = -diff
		}
		if diff > eps*maxFloat(1.0, expect) {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%.10f actual:%.10f\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
