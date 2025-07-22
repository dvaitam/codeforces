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

type testCaseC struct {
	input    string
	expected int
}

func solveC(n int) int {
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		best := 1<<31 - 1
		x := i
		for x > 0 {
			d := x % 10
			x /= 10
			if d > 0 {
				if dp[i-d] < best {
					best = dp[i-d]
				}
			}
		}
		dp[i] = best + 1
	}
	return dp[n]
}

func genCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(1000) + 1
	expected := solveC(n)
	return testCaseC{input: fmt.Sprintf("%d\n", n), expected: expected}
}

func runCaseC(bin string, tc testCaseC) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
