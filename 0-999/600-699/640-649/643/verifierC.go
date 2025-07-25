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

type testCase struct {
	input  string
	n, k   int
	t      []int
	expect float64
}

func cost(S, P1, P2 []float64, l, r int) float64 {
	return P1[r] - P1[l-1] - S[l-1]*(P2[r]-P2[l-1])
}

func solve(n, k int, t []int) float64 {
	S := make([]float64, n+1)
	P1 := make([]float64, n+1)
	P2 := make([]float64, n+1)
	for i := 1; i <= n; i++ {
		S[i] = S[i-1] + float64(t[i-1])
		P1[i] = P1[i-1] + S[i]/float64(t[i-1])
		P2[i] = P2[i-1] + 1.0/float64(t[i-1])
	}
	inf := math.Inf(1)
	dp := make([][]float64, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]float64, n+1)
		for j := 0; j <= n; j++ {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0
	for c := 1; c <= k; c++ {
		for i := 1; i <= n; i++ {
			for j := 0; j < i; j++ {
				v := dp[c-1][j] + cost(S, P1, P2, j+1, i)
				if v < dp[c][i] {
					dp[c][i] = v
				}
			}
		}
	}
	return dp[k][n]
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	k := rng.Intn(n) + 1
	t := make([]int, n)
	for i := 0; i < n; i++ {
		t[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", t[i])
	}
	sb.WriteByte('\n')
	expect := solve(n, k, t)
	return testCase{sb.String(), n, k, t, expect}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseFloat(out string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(out), 64)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		val, err := parseFloat(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse float: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		expect := c.expect
		if math.Abs(val-expect) > 1e-4*math.Max(1, math.Abs(expect)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", i+1, expect, val, c.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
