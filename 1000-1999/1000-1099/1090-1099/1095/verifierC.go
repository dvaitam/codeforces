package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	n, k  int64
	input string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1095 + 2))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Int63n(1e9) + 1
		limit := n
		if limit > 200000 {
			limit = 200000
		}
		k := rng.Int63n(limit) + 1
		input := fmt.Sprintf("%d %d\n", n, k)
		cases[i] = Case{n: n, k: k, input: input}
	}
	return cases
}

func isPow2(x int64) bool {
	return x > 0 && (x&(x-1)) == 0
}

func popcount(x int64) int64 {
	var c int64
	for x > 0 {
		c += x & 1
		x >>= 1
	}
	return c
}

func validateCase(c Case, output string) error {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	firstLine := strings.TrimSpace(strings.ToUpper(lines[0]))

	// Determine if answer should be YES or NO
	minK := popcount(c.n)
	canYes := c.k >= minK && c.k <= c.n

	if firstLine == "NO" {
		if canYes {
			return fmt.Errorf("candidate said NO but answer should be YES (n=%d k=%d minK=%d)", c.n, c.k, minK)
		}
		return nil
	}
	if firstLine != "YES" {
		return fmt.Errorf("expected YES or NO, got %q", firstLine)
	}
	if !canYes {
		return fmt.Errorf("candidate said YES but answer should be NO (n=%d k=%d)", c.n, c.k)
	}

	// Parse the values from remaining lines (could be on same line or next line)
	var valStr string
	if len(lines) > 1 {
		valStr = strings.Join(lines[1:], " ")
	}
	fields := strings.Fields(valStr)
	if int64(len(fields)) != c.k {
		return fmt.Errorf("expected %d values, got %d", c.k, len(fields))
	}
	var sum int64
	for _, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value %q", f)
		}
		if !isPow2(v) {
			return fmt.Errorf("value %d is not a power of 2", v)
		}
		sum += v
	}
	if sum != c.n {
		return fmt.Errorf("sum %d != n %d", sum, c.n)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := genCases()
	for i, c := range cases {
		got, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := validateCase(c, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, c.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
