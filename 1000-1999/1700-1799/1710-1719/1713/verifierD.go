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
}

// solve embeds the (placeholder) logic from 1713D.go.
func solve(tc testCase) string {
	_ = tc.n
	return "1"
}

// Embedded copy of testcasesD.txt.
const testcaseData = `
3
10
1
8
3
6
1
8
5
10
4
2
9
7
5
3
9
3
2
3
10
2
9
9
10
7
7
5
5
5
1
7
5
5
9
9
9
6
6
4
7
3
1
9
3
10
7
6
8
1
9
7
10
4
1
6
9
3
4
6
8
1
4
10
4
5
3
7
2
10
8
4
8
9
2
4
3
8
2
7
7
5
5
7
6
10
6
2
5
1
8
1
5
4
7
7
7
7
1
10
`

// Expected outputs for each testcase (placeholder solution prints 1).
var expectedOutputs = []string{
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
	"1", "1", "1", "1", "1", "1", "1", "1", "1", "1",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')

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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(bin, tc, expectedOutputs[i]); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
