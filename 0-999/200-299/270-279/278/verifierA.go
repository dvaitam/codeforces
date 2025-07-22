package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// runBinary executes the given binary with provided input and returns its output or error.
func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// expectedA computes the correct answer for problem A given the input.
func expectedA(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return "", err
	}
	d := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &d[i]); err != nil {
			return "", err
		}
	}
	var s, t int
	if _, err := fmt.Fscan(in, &s, &t); err != nil {
		return "", err
	}
	if s == t {
		return "0\n", nil
	}
	if s > t {
		s, t = t, s
	}
	clockwise := 0
	for i := s - 1; i < t-1; i++ {
		clockwise += d[i]
	}
	total := 0
	for _, x := range d {
		total += x
	}
	counter := total - clockwise
	if clockwise < counter {
		return fmt.Sprintf("%d\n", clockwise), nil
	}
	return fmt.Sprintf("%d\n", counter), nil
}

// generateCase produces a deterministic random test case.
func generateCase(rng *rand.Rand) string {
	n := rng.Intn(98) + 3 // 3..100
	d := make([]int, n)
	for i := 0; i < n; i++ {
		d[i] = rng.Intn(100) + 1
	}
	s := rng.Intn(n) + 1
	t := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, x := range d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", x)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))

	var cases []string
	// some fixed edge cases
	cases = append(cases, "3\n1 1 1\n1 3\n")
	cases = append(cases, "3\n100 100 100\n2 2\n")
	cases = append(cases, "4\n1 2 3 4\n1 4\n")
	cases = append(cases, "5\n1 2 3 4 5\n5 1\n")

	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		expect, err := expectedA(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected: %sgot: %s", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
