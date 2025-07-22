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

type testCase struct {
	x, k, p int
}

func (tc testCase) Input() string {
	return fmt.Sprintf("%d %d %d\n", tc.x, tc.k, tc.p)
}

func generateCase(rng *rand.Rand) testCase {
	x := rng.Intn(100) + 1
	k := rng.Intn(5) + 1
	p := rng.Intn(101)
	return testCase{x, k, p}
}

func expected(tc testCase) float64 {
	pd := float64(tc.p) / 100.0
	pp := float64(100-tc.p) / 100.0
	const MAX = 310
	dp := make([][]float64, 2)
	dp[0] = make([]float64, MAX)
	dp[1] = make([]float64, MAX)
	dp[0][0] = 1.0
	next := 1
	ans := 0.0
	for i := 0; i < tc.k; i++ {
		prev := 1 - next
		for j := 0; j < MAX; j++ {
			dp[next][j] = 0.0
		}
		for j := 0; j < 300; j++ {
			dp[next][j+1] += dp[prev][j] * pp
		}
		for j := 0; j < 300; j += 2 {
			ans += dp[prev][j] * pd
			dp[next][j/2] += dp[prev][j] * pd
		}
		next = prev
	}
	prev := 1 - next
	for i := 0; i < 300; i++ {
		r := i + tc.x
		for r%2 == 0 {
			r /= 2
			ans += dp[prev][i]
		}
	}
	return ans
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	input := tc.Input()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expected(tc)
	diff := got - expect
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-6*maxFloat64(expect, 1.0) && diff > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", expect, got)
	}
	return nil
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{{1, 1, 0}, {5, 3, 50}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
