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
	input  string
	output string
}

func solveCase(n, k int64) int64 {
	mult := (n + k - 1) / k
	sum := mult * k
	ans := (sum + n - 1) / n
	return ans
}

func buildCase(n, k int64) testCase {
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	output := fmt.Sprintf("%d\n", solveCase(n, k))
	return testCase{input: input, output: output}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Int63n(1000) + 1
	k := rng.Int63n(1000) + 1
	return buildCase(n, k)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// some deterministic cases
	cases = append(cases, buildCase(1, 5))
	cases = append(cases, buildCase(4, 3))
	cases = append(cases, buildCase(8, 8))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
