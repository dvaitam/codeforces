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
4
18
2
8
3
15
14
15
20
12
6
3
15
0
12
13
19
0
14
8
7
18
3
10
0
0
0
20
17
0
12
6
13
0
16
7
14
15
17
7
11
7
7
14
9
0
13
17
20
3
5
20
9
3
10
16
13
16
6
9
9
18
15
16
12
18
1
15
7
12
13
5
11
17
11
2
14
16
3
5
16
12
11
15
0
15
1
9
19
18
18
12
20
5
5
16
7
0
6
17
`

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesCData), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
}

func solve(n int) uint64 {
	var r uint64 = 1
	for i := 2; i <= n; i++ {
		r *= uint64(i)
	}
	return r
}

func runCandidate(bin string, n int) (uint64, error) {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	valStr := strings.TrimSpace(out.String())
	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", valStr, err)
	}
	return val, nil
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
		expect := solve(tc.n)
		got, err := runCandidate(bin, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
