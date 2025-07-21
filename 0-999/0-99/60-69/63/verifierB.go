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

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

func compute(n, k int, ranks []int) int {
	counts := make([]int, k+2)
	for _, r := range ranks {
		if r >= 1 && r <= k {
			counts[r]++
		}
	}
	sessions := 0
	for {
		done := true
		for r := 1; r < k; r++ {
			if counts[r] > 0 {
				done = false
				break
			}
		}
		if done {
			break
		}
		sessions++
		newC := make([]int, k+2)
		copy(newC, counts)
		for r := 1; r < k; r++ {
			if counts[r] > 0 {
				newC[r]--
				newC[r+1]++
			}
		}
		counts = newC
	}
	return sessions
}

func generateCases() []testCase {
	rand.Seed(2)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(9) + 2
		ranks := make([]int, n)
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			ranks[j] = rand.Intn(k) + 1
			fmt.Fprintf(&buf, "%d", ranks[j])
			if j+1 < n {
				fmt.Fprint(&buf, " ")
			}
		}
		buf.WriteByte('\n')
		expected := fmt.Sprintf("%d", compute(n, k, ranks))
		cases[i] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB.go <binary>")
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
