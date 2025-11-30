package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	s string
}

// Embedded logic from 1704F.go.
func solve(tc testCase) string {
	n := tc.n
	s := tc.s
	r, b := 0, 0
	for _, ch := range s {
		if ch == 'R' {
			r++
		} else if ch == 'B' {
			b++
		}
	}
	if r > b {
		return "Alice"
	}
	if b > r {
		return "Bob"
	}
	isPal := true
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			isPal = false
			break
		}
	}
	if isPal {
		return "Bob"
	}
	return "Alice"
}

// Embedded copy of testcasesF.txt (each line: n string).
const testcaseData = `
8 BBRRBBBB
7 BBBRRBB
4 RRBR
3 BRR
3 BRR
5 BBRRR
5 BBBRR
5 BRRRR
5 RBRBB
3 BBR
6 BRRRBR
5 BRBBB
1 R
4 BRRR
1 B
5 RRRRR
3 RBR
3 BBR
8 BBRRRRBB
4 BRBR
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields got %d", i+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n, s: fields[1]})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d\n%s\n", tc.n, tc.s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
