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
	input       string
	expectedOps int
	expectedSeg int
}

func solveCase(n int, s string) (int, int) {
	k := n / 2
	ops := 0
	pred := make([]byte, 0, k)
	for i := 0; i < k; i++ {
		a := s[2*i]
		b := s[2*i+1]
		if a != b {
			ops++
		} else {
			pred = append(pred, a)
		}
	}
	seg := 1
	if len(pred) > 0 {
		seg = 1
		for i := 1; i < len(pred); i++ {
			if pred[i] != pred[i-1] {
				seg++
			}
		}
	}
	return ops, seg
}

func buildCase(s string) testCase {
	var sb strings.Builder
	n := len(s)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	ops, seg := solveCase(n, s)
	return testCase{input: sb.String(), expectedOps: ops, expectedSeg: seg}
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
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 numbers got %d", len(fields))
	}
	var gotOps, gotSeg int
	if _, err := fmt.Sscan(fields[0], &gotOps); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if _, err := fmt.Sscan(fields[1], &gotSeg); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if gotOps != tc.expectedOps || gotSeg != tc.expectedSeg {
		return fmt.Errorf("expected %d %d got %d %d", tc.expectedOps, tc.expectedSeg, gotOps, gotSeg)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
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
