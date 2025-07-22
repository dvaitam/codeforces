package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input string
	n, T  int
	p, t  []int
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n, T int, p, t []int) float64 {
	dp := make([][]float64, 2)
	dp[0] = make([]float64, T+1)
	dp[1] = make([]float64, T+1)
	cur := 0
	for idx := n - 1; idx >= 0; idx-- {
		cur ^= 1
		prev := cur ^ 1
		prob := float64(p[idx]) / 100.0
		P := math.Pow(1-prob, float64(t[idx]-1))
		sum := 0.0
		dp[cur][0] = 0
		for j := 1; j <= T; j++ {
			sum *= 1 - prob
			sum += (dp[prev][j-1] + 1) * prob
			if j >= t[idx] {
				sum += (dp[prev][j-t[idx]] + 1) * P * (1 - prob)
			}
			if j > t[idx] {
				sum -= (dp[prev][j-t[idx]-1] + 1) * P * (1 - prob)
			}
			dp[cur][j] = sum
		}
	}
	return dp[cur][T]
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	T := rng.Intn(20) + 1
	p := make([]int, n)
	t := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Intn(101)
		t[i] = rng.Intn(T) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, T))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[i], t[i]))
	}
	return testCase{input: sb.String(), n: n, T: T, p: p, t: t}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 101)
	cases = append(cases, testCase{input: "1 1\n100 1\n", n: 1, T: 1, p: []int{100}, t: []int{1}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		outStr, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var outVal float64
		fmt.Sscanf(outStr, "%f", &outVal)
		exp := expected(tc.n, tc.T, tc.p, tc.t)
		if math.Abs(outVal-exp) > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %.6f\nfound: %s\n", i+1, tc.input, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
