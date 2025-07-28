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

type testCaseB struct {
	n int64
	m int64
	k int64
}

func expectedB(tc testCaseB) string {
	n, m, k := tc.n, tc.m, tc.k
	maxEdges := n * (n - 1) / 2
	if m < n-1 || m > maxEdges {
		return "NO"
	}
	if n == 1 {
		if m == 0 && k > 1 {
			return "YES"
		}
		return "NO"
	}
	if n == 2 {
		if m == 1 && k > 2 {
			return "YES"
		}
		return "NO"
	}
	var minDiameter int64
	if m == maxEdges {
		minDiameter = 1
	} else {
		minDiameter = 2
	}
	if k-1 > minDiameter {
		return "YES"
	}
	return "NO"
}

func buildInputB(tc testCaseB) string {
	return fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.m, tc.k)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomCaseB(rng *rand.Rand) testCaseB {
	n := rng.Int63n(1_000_000_000) + 1
	maxEdges := n * (n - 1) / 2
	if maxEdges > 1_000_000_000 {
		maxEdges = 1_000_000_000
	}
	m := rng.Int63n(maxEdges + 1)
	k := rng.Int63n(1_000_000_000) + 1
	return testCaseB{n: n, m: m, k: k}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCaseB
	// deterministic cases
	cases = append(cases, testCaseB{n: 1, m: 0, k: 2})
	cases = append(cases, testCaseB{n: 1, m: 1, k: 2})
	cases = append(cases, testCaseB{n: 2, m: 1, k: 3})
	cases = append(cases, testCaseB{n: 2, m: 1, k: 2})
	cases = append(cases, testCaseB{n: 3, m: 2, k: 3})
	cases = append(cases, testCaseB{n: 3, m: 3, k: 4})

	for len(cases) < 110 {
		cases = append(cases, randomCaseB(rng))
	}

	for i, tc := range cases {
		input := buildInputB(tc)
		expect := expectedB(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(out)) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
