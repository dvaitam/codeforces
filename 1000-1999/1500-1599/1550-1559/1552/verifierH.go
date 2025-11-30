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
	line string
}

// Embedded testcases from testcasesH.txt.
const testcaseData = `
1 9
6 1
2 7
6 8
10 1
6 6
3 9
2 9
4 3
6 4
9 2
8 10
4 4
6 6
5 9
7 5
1 2
6 6
8 8
9 9
2 8
9 9
5 8
10 6
2 4
8 2
2 7
3 5
2 2
3 5
2 3
8 3
4 9
1 5
7 9
3 8
6 7
8 1
5 4
2 1
4 9
6 7
10 6
6 8
10 9
9 8
8 5
10 3
5 3
5 4
8 6
5 1
3 6
5 8
3 10
3 9
5 3
3 3
3 4
2 1
6 10
3 8
5 6
3 4
9 1
6 6
3 6
8 9
3 9
5 2
1 7
8 3
8 6
3 1
9 9
6 5
8 5
10 3
10 8
6 2
8 1
3 9
2 5
1 7
10 4
8 6
2 6
1 10
4 10
9 4
5 4
1 2
9 8
2 9
7 4
8 6
10 10
9 1
6 7
3 7
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("case %d: expected two integers, got %d", idx+1, len(fields))
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return nil, fmt.Errorf("case %d bad value: %v", idx+1, err)
		}
		if _, err := strconv.Atoi(fields[1]); err != nil {
			return nil, fmt.Errorf("case %d bad value: %v", idx+1, err)
		}
		cases = append(cases, testCase{line: line})
	}
	return cases, nil
}

// solve mirrors 1552H.go placeholder.
func solve(_ testCase) string {
	return "0"
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := tc.line + "\n"
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
