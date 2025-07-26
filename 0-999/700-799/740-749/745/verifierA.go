package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func countShifts(s string) int {
	n := len(s)
	seen := make(map[string]struct{})
	for i := 0; i < n; i++ {
		t := s[n-i:] + s[:n-i]
		seen[t] = struct{}{}
	}
	return len(seen)
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

type testCase struct {
	input    string
	expected string
}

func genCase(rng *rand.Rand) testCase {
	l := rng.Intn(50) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	input := fmt.Sprintf("%s\n", s)
	expected := fmt.Sprintf("%d", countShifts(s))
	return testCase{input: input, expected: expected}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(745))

	cases := []testCase{
		{input: "abcd\n", expected: fmt.Sprintf("%d", countShifts("abcd"))},
		{input: "bbb\n", expected: fmt.Sprintf("%d", countShifts("bbb"))},
		{input: "yzyz\n", expected: fmt.Sprintf("%d", countShifts("yzyz"))},
		{input: "a\n", expected: fmt.Sprintf("%d", countShifts("a"))},
		{input: "aba\n", expected: fmt.Sprintf("%d", countShifts("aba"))},
	}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, tc.expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
