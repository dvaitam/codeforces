package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
10 21
5 0
5 0
8 10
7 1
10 18
7 3
10 22
10 17
10 9
6 1
5 0
10 10
6 2
10 15
7 2
6 0
9 16
7 5
9 0
6 3
6 1
9 2
6 1
8 9
10 6
6 1
6 3
7 2
10 15
5 0
9 11
8 7
7 0
6 1
10 15
10 15
9 10
10 2
7 2
10 1
10 15
10 9
8 2
5 0
5 0
6 4
5 0
7 3
8 8
7 2
9 1
8 3
9 9
9 5
8 9
10 16
8 7
9 13
10 15
10 22
6 0
8 3
8 9
8 1
7 2
7 0
8 8
9 11
10 5
10 16
10 19
6 2
9 13
6 0
7 2
10 6
6 0
6 1
8 0
10 7
9 15
5 0
6 2
8 7
6 2
8 8
8 5
6 1
6 0
6 1
8 1
7 2
10 5
10 9
9 2
10 12
5 0
10 17
9 4
8 9
10 9
9 2
5 0
10 14
6 1
5 0
10 16
`

type testCase struct {
	n int
	p int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesCData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 integers", idx+1)
		}
		n, _ := strconv.Atoi(fields[0])
		p, _ := strconv.Atoi(fields[1])
		cases = append(cases, testCase{n: n, p: p})
	}
	return cases, nil
}

// solve mirrors 402C.go for a single test case.
func solve(n, p int) string {
	var sb strings.Builder
	// distance 1 edges
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		fmt.Fprintf(&sb, "%d %d\n", i+1, j+1)
	}
	// distance 2 edges
	for i := 0; i < n; i++ {
		j := (i + 2) % n
		fmt.Fprintf(&sb, "%d %d\n", i+1, j+1)
	}
	// extra edges with d>=3
	need := p
	for d := 3; d < n && need > 0; d++ {
		for i := 0; i < n && need > 0; i++ {
			j := (i + d) % n
			if i < j {
				fmt.Fprintf(&sb, "%d %d\n", i+1, j+1)
				need--
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.p)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.n, tc.p)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
