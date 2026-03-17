package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"math/rand"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + r.Intn(26))
	}
	return string(b)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(100) + 1
	s := randString(r, n)
	return fmt.Sprintf("%s\n", s)
}

func validateCase(input, output string) error {
	s := strings.TrimSpace(input)
	L := len(s)

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return fmt.Errorf("empty output")
	}

	var a, b int
	if _, err := fmt.Sscan(lines[0], &a, &b); err != nil {
		return fmt.Errorf("invalid first line: %v", err)
	}

	// Check minimum rows: a <= 5, b <= 20
	// Minimum a is ceil(L / 20), but at least 1
	minA := (L + 19) / 20
	if minA < 1 {
		minA = 1
	}
	if a != minA {
		return fmt.Errorf("expected %d rows, got %d", minA, a)
	}

	// Minimum b for given a: ceil(L / a), accounting for asterisks
	// Total cells = a * b, asterisks = a*b - L
	// asterisks per row differ by at most 1
	// b = ceil(L / a) ... but with asterisk constraint
	minB := (L + a - 1) / a
	totalStars := a*minB - L
	if totalStars > a {
		// This shouldn't happen with correct minA
		return fmt.Errorf("unexpected star count")
	}
	if b != minB {
		return fmt.Errorf("expected %d columns, got %d", minB, b)
	}

	if len(lines) != a+1 {
		return fmt.Errorf("expected %d table rows, got %d", a, len(lines)-1)
	}

	// Check each row has exactly b characters
	// Count asterisks per row
	starCounts := make([]int, a)
	var extracted strings.Builder
	for i := 0; i < a; i++ {
		row := lines[i+1]
		if len(row) != b {
			return fmt.Errorf("row %d: expected %d chars, got %d", i+1, b, len(row))
		}
		for _, ch := range row {
			if ch == '*' {
				starCounts[i]++
			} else {
				extracted.WriteRune(ch)
			}
		}
	}

	// Check extracted handle
	if extracted.String() != s {
		return fmt.Errorf("extracted handle %q != expected %q", extracted.String(), s)
	}

	// Check asterisk counts differ by at most 1
	minStars, maxStars := starCounts[0], starCounts[0]
	for _, c := range starCounts {
		if c < minStars {
			minStars = c
		}
		if c > maxStars {
			maxStars = c
		}
	}
	if maxStars-minStars > 1 {
		return fmt.Errorf("asterisk counts vary by more than 1: min=%d max=%d", minStars, maxStars)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	cases := []string{
		"a\n",
		"ab\n",
		"abcde\n",
	}
	for i := 0; i < 97; i++ {
		cases = append(cases, genCase(r))
	}
	for idx, input := range cases {
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if err := validateCase(input, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
