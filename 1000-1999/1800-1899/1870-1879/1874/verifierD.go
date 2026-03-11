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

func runProg(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// solve computes the minimum expected number of days for given n, m.
// Expected days = n + 2 * sum_{i=2}^{n} (a_1+...+a_{i-1})/a_i
// Subject to: a_i >= 1 (positive integers), sum(a_i) <= m.
// We minimize sum_{i=2}^{n} S_{i-1}/a_i where S_{i-1} = a_1+...+a_{i-1}.
func solve(n, m int) float64 {
	if n == 1 {
		return 1.0
	}

	const INF = 1e18

	// DP: f[i][j] = minimum sum of S_{prev}/a_i for roads 1..i with total beauty j
	// i ranges from 1 to n, j from 1 to m
	// f[1][j] = 0 for j = 1..m (road 1 has beauty j, no cost term for road 1)
	// f[i][j] = min over b=1..j-(i-1) of f[i-1][j-b] + (j-b)/b  (beauty of road i = b, cumulative before = j-b)
	// But we also need j >= i (at least 1 beauty per road)

	f := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]float64, m+1)
		for j := 0; j <= m; j++ {
			f[i][j] = INF
		}
	}

	// Base: f[1][j] = 0 for j >= 1
	for j := 1; j <= m; j++ {
		f[1][j] = 0
	}

	// Fill DP
	for i := 2; i <= n; i++ {
		for j := i; j <= m; j++ {
			// Road i has beauty b, previous roads use j-b total beauty
			// previous roads: i-1 roads with total beauty j-b, need j-b >= i-1
			for b := 1; b <= j-(i-1); b++ {
				prev := j - b
				if f[i-1][prev] >= INF {
					continue
				}
				cost := f[i-1][prev] + float64(prev)/float64(b)
				if cost < f[i][j] {
					f[i][j] = cost
				}
			}
		}
	}

	ans := INF
	for j := n; j <= m; j++ {
		if f[n][j] < ans {
			ans = f[n][j]
		}
	}

	return float64(n) + 2*ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Test with known sample first
	sampleInput := []byte("3 8\n")
	sampleExpected := 5.2
	got, err := runProg(cand, sampleInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed on sample: %v\n", err)
		os.Exit(1)
	}
	gotVal, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output on sample: %v\noutput: %s\n", err, got)
		os.Exit(1)
	}
	if math.Abs(gotVal-sampleExpected)/math.Max(1.0, math.Abs(sampleExpected)) > 1e-6 {
		fmt.Fprintf(os.Stderr, "sample test failed: expected %.6f got %.6f\n", sampleExpected, gotVal)
		os.Exit(1)
	}

	// Random tests using our correct DP solver
	for i := 1; i <= 100; i++ {
		n := rng.Intn(15) + 1
		m := n + rng.Intn(30)
		if m > 50 {
			m = 50 // keep small for O(n*m^2) DP
		}
		if m < n {
			m = n
		}
		input := []byte(fmt.Sprintf("%d %d\n", n, m))
		wantVal := solve(n, m)

		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		gotVal, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput: %s\n", i, err, got)
			os.Exit(1)
		}
		denom := math.Max(1.0, math.Abs(wantVal))
		if math.Abs(wantVal-gotVal)/denom > 1e-6 {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput: %sexpected: %.12f\ngot: %.12f\n", i, string(input), wantVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
