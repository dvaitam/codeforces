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
	input    string
	expected string
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(xs, ys []int, k int) int {
	n := len(xs)
	for i := 0; i < n; i++ {
		ok := true
		for j := 0; j < n; j++ {
			if abs(xs[i]-xs[j])+abs(ys[i]-ys[j]) > k {
				ok = false
				break
			}
		}
		if ok {
			return 1
		}
	}
	return -1
}

func buildCase(xs, ys []int, k int) testCase {
	n := len(xs)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
	}
	ans := solveCase(xs, ys, k)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d", ans)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	k := rng.Intn(20)
	xs := make([]int, n)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = rng.Intn(21)
		ys[i] = rng.Intn(21)
	}
	return buildCase(xs, ys, k)
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
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
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
