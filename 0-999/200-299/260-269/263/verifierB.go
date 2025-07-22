package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func generateCases() []testCase {
	rand.Seed(2)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		k := rand.Intn(50) + 1
		vals := make([]int, 0, n)
		used := make(map[int]bool)
		for len(vals) < n {
			v := rand.Intn(1000000) + 1
			if !used[v] {
				used[v] = true
				vals = append(vals, v)
			}
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, k)
		for j, v := range vals {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
		sorted := append([]int(nil), vals...)
		sort.Ints(sorted)
		var exp string
		if k > n {
			exp = "-1"
		} else {
			exp = fmt.Sprintf("%d 0", sorted[n-k])
		}
		cases[i] = testCase{input: buf.String(), expected: exp}
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
