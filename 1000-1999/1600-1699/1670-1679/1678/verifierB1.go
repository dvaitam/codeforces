package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input    string
	expected int
}

func solveCase(n int, s string) int {
	ops := 0
	for i := 0; i < n; i += 2 {
		if s[i] != s[i+1] {
			ops++
		}
	}
	return ops
}

func buildCase(s string) testCase {
	var sb strings.Builder
	n := len(s)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	return testCase{input: sb.String(), expected: solveCase(n, s)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100)*2 + 2
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return buildCase(string(b))
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
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	var val int
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	var cases []testCase

	// edge cases
	cases = append(cases, buildCase("00"))
	cases = append(cases, buildCase("01"))
	cases = append(cases, buildCase("10"))
	cases = append(cases, buildCase("11"))

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
