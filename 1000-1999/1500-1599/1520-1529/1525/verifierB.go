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
	n   int
	arr []int
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
5 3 4 2 1 5
10 3 6 5 4 1 2 9 7 10 8
9 7 4 5 9 2 6 3 8 1
3 2 3 1
9 8 3 9 2 5 1 6 7 4
6 1 6 4 5 2 3
3 3 1 2
5 4 2 5 1 3
9 1 2 6 8 5 7 3 4 9
3 3 1 2
9 6 4 2 1 9 5 8 3 7
5 1 3 2 4 5
3 3 1 2
7 2 1 3 4 7 5 6
5 4 3 1 2 5
6 1 2 6 3 5 4
10 7 8 3 2 10 6 4 1 9 5
9 9 4 6 2 5 3 7 8 1
10 4 9 2 3 8 5 1 10 7 6
6 3 4 1 5 2 6
3 2 1 3
4 3 4 2 1
7 5 6 4 7 1 3 2
4 3 4 1 2
5 1 4 2 5 3
10 5 9 10 6 7 4 3 1 2 8
4 2 1 4 3
9 5 6 3 7 2 8 9 4 1
10 2 3 1 8 5 10 6 4 7 9
9 4 5 8 2 3 9 1 7 6
7 2 6 4 5 3 7 1
9 6 4 8 2 7 9 1 3 5
5 3 4 2 1 5
8 3 1 6 7 4 5 8 2
9 6 7 2 4 9 3 1 8 5
5 5 1 4 3 2
8 5 8 6 4 1 3 2 7
10 10 3 5 8 7 1 6 2 4 9
6 1 4 2 5 6 3
8 6 7 1 4 5 8 3 2
8 3 5 6 2 7 8 4 1
4 2 3 1 4
4 1 4 2 3
4 3 1 2 4
3 3 1 2
3 3 2 1
4 2 3 4 1
9 7 1 9 5 6 8 4 2 3
9 2 9 1 6 4 3 8 5 7
8 6 4 2 5 3 7 8 1
8 5 2 6 4 3 8 1 7
6 1 2 3 4 6 5
7 6 4 1 3 7 5 2
8 6 8 5 4 1 3 7 2
6 5 6 2 1 3 4
8 7 6 4 8 1 3 2 5
3 3 2 1
4 2 1 4 3
4 2 4 3 1
7 2 6 1 5 4 3 7
8 4 3 1 7 8 6 5 2
8 2 4 7 3 8 6 1 5
5 3 4 2 1 5
5 2 1 5 4 3
3 2 1 3
7 6 2 3 5 4 1 7
10 7 10 5 6 2 3 4 1 8 9
9 4 8 7 2 5 3 6 9 1
5 1 4 2 5 3
5 3 2 4 1 5
3 3 2 1
10 1 5 3 2 7 8 9 6 10 4
8 4 7 3 8 1 2 6 5
5 1 3 4 2 5
7 1 4 6 2 3 5 7
4 4 2 3 1
5 4 1 5 3 2
3 1 3 2
8 6 4 3 1 8 2 5 7
7 2 5 4 7 3 1 6
10 5 6 1 8 3 2 10 9 7 4
6 3 2 6 5 4 1
7 3 7 6 1 5 4 2
3 2 3 1
7 4 1 5 6 3 7 2
7 6 2 1 4 7 5 3
10 10 6 1 3 9 8 5 4 2 7
4 4 2 3 1
5 5 2 4 3 1
9 8 2 4 1 3 5 7 6 9
10 8 2 1 4 9 10 3 7 5 6
9 8 7 9 2 5 6 3 1 4
7 3 1 5 7 4 2 6
9 8 4 2 1 6 5 3 9 7
7 7 5 3 2 4 1 6
7 4 3 6 5 1 2 7
9 8 7 1 6 3 4 2 5 9
3 3 1 2
9 9 3 5 7 4 8 1 6 2
4 1 2 3 4
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1525B.go.
func solve(tc testCase) int {
	n := tc.n
	arr := tc.arr
	sorted := true
	for i, v := range arr {
		if v != i+1 {
			sorted = false
			break
		}
	}
	if sorted {
		return 0
	}
	if arr[0] == 1 || arr[n-1] == n {
		return 1
	}
	if arr[0] == n && arr[n-1] == 1 {
		return 3
	}
	return 2
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if got != strconv.Itoa(expect) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
