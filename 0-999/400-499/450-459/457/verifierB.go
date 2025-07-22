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

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var m, n int
	fmt.Fscan(rdr, &m, &n)
	sumA, sumB := uint64(0), uint64(0)
	maxA, maxB := uint64(0), uint64(0)
	for i := 0; i < m; i++ {
		var x uint64
		fmt.Fscan(rdr, &x)
		sumA += x
		if x > maxA {
			maxA = x
		}
	}
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(rdr, &x)
		sumB += x
		if x > maxB {
			maxB = x
		}
	}
	cost1 := (sumA - maxA) + sumB
	cost2 := sumA + (sumB - maxB)
	cost3 := sumB * uint64(m)
	cost4 := sumA * uint64(n)
	ans := cost1
	if cost2 < ans {
		ans = cost2
	}
	if cost3 < ans {
		ans = cost3
	}
	if cost4 < ans {
		ans = cost4
	}
	return fmt.Sprintf("%d", ans)
}

func generateCases() []testCase {
	rand.Seed(2)
	cases := []testCase{}
	fixed := []string{
		"1 1\n5\n7\n",
		"2 2\n1 2\n3 4\n",
		"3 1\n10 20 30\n5\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		m := rand.Intn(4) + 1
		n := rand.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", m, n)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d", rand.Intn(100)+1)
			if i+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d", rand.Intn(100)+1)
			if i+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
