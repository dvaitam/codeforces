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

// Embedded testcases from testcasesF.txt.
const testcaseData = `
100
2
10
6
1
10
3
5
3
6
8
4
7
6
9
1
6
6
10
8
4
8
7
10
10
6
7
7
8
4
3
8
1
4
1
8
1
3
2
3
9
4
2
2
7
4
9
1
3
7
3
2
7
6
9
4
8
1
2
2
4
9
2
2
9
2
10
2
1
3
5
8
10
9
3
5
5
10
10
10
4
6
5
3
8
3
6
4
4
3
10
7
8
7
1
10
3
10
1
5
2
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		return nil, fmt.Errorf("no test data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	if len(lines)-1 != t {
		return nil, fmt.Errorf("expected %d test cases, got %d", t, len(lines)-1)
	}
	cases := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		n, err := strconv.Atoi(strings.TrimSpace(lines[i]))
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
}

// solve mirrors 1526F.go (placeholder interactive solver).
func solve(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("!")
	for i := 1; i <= tc.n; i++ {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(i))
	}
	return sb.String()
}

func runCandidate(bin string, tc testCase) (string, error) {
	inp := fmt.Sprintf("1\n%d\n", tc.n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(inp)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
