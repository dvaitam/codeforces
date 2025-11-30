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
	a []int
	b []int
}

// solve embeds the logic from 1736A.go.
func solve(tc testCase) string {
	cntA, cntB := 0, 0
	for _, v := range tc.a {
		if v == 1 {
			cntA++
		}
	}
	for _, v := range tc.b {
		if v == 1 {
			cntB++
		}
	}
	mism := 0
	for i := 0; i < tc.n; i++ {
		if tc.a[i] != tc.b[i] {
			mism++
		}
	}
	diff := cntA - cntB
	if diff < 0 {
		diff = -diff
	}
	ans := mism
	if diff+1 < ans {
		ans = diff + 1
	}
	return strconv.Itoa(ans)
}

// Embedded copy of testcasesA.txt.
const testcaseData = `
7 1 0 1 1 1 1 1 1 0 0 1 0 0 1
9 0 1 0 0 1 1 0 1 1 1 0 1 1 1 0 0 0 1
1 1 1
4 1 0 0 0 0 0 1 0
2 1 1 0 1
9 1 0 1 0 1 1 0 1 1 0 1 0 0 0 0 1 1 0
2 0 0 0 0
9 1 1 0 0 1 1 1 1 1 0 1 0 1 1 0 0 0 1
2 0 1 0 1
7 0 0 0 0 0 0 0 0 0 0 1 0 1 0
1 0 0
3 0 1 0 0 0 1
10 0 1 0 0 0 1 1 1 0 0 1 0 0 1 0 1 1 1 0 0
1 0 0
6 1 0 1 0 0 1 1 1 1 1 1 0
9 0 1 0 1 0 1 0 0 1 1 1 1 0 1 1 1 0 0
10 0 1 0 0 0 1 1 1 0 1 1 0 0 1 0 1 0 1 1 0
1 1 1
5 1 0 1 0 0 0 0 1 1 1
1 0 0
1 0 0
2 0 1 1 0
2 1 1 0 0
5 1 0 1 0 1 0 0 1 0 0
1 0 1
9 1 1 0 1 1 1 1 0 0 1 1 0 0 0 1 1 1 1
2 1 0 0 1
3 0 1 1 1 0 1
2 1 0 0 1
3 0 1 1 1 1 1
2 0 1 1 1
6 0 1 0 1 1 0 1 1 0 0 1 0
2 0 0 0 0
7 0 0 1 1 1 1 0 1 0 1 0 1 0 1
4 1 0 0 0 1 0 0 1
7 1 0 0 0 0 0 0 1 1 1 0 0 1 0
10 1 1 1 1 1 1 0 0 0 1 1 1 0 0 1 0 1 1 1 0
2 0 0 0 0
1 1 0
8 1 0 0 1 1 1 0 0 0 1 0 1 1 1 1 0
1 0 0
5 1 1 0 1 1 1 1 1 1 0
3 0 0 1 1 0 0
8 1 0 1 0 0 1 1 0 1 1 1 0 0 1 0 0
1 0 1
2 1 0 1 1
2 0 1 1 0
10 1 0 1 1 0 0 0 1 1 1 0 1 1 0 0 0 1 1 0 0
4 0 1 1 0 0 1 1 0
2 1 0 0 1
3 0 1 1 1 0 1
6 1 0 1 0 0 1 0 1 1 0 0 0
6 0 0 0 0 0 0 0 1 1 0 0 0
3 1 0 1 1 1 0
1 0 1
6 1 1 1 1 0 0 0 1 0 1 0 0
4 1 0 0 0 1 1 1 0
3 0 1 0 0 0 1
5 1 0 1 1 1 0 0 1 1 0
7 1 0 1 1 1 0 1 1 1 0 0 0 1 0
9 0 1 1 0 1 1 0 0 1 1 0 1 1 0 1 0 1 0
8 1 0 0 1 0 1 0 1 1 0 1 0 1 1 0 0
5 0 0 1 1 1 1 0 1 1 1
4 0 0 0 1 1 1 0 0
3 0 1 1 1 0 0
9 0 1 1 1 1 1 1 0 0 0 0 1 1 1 0 1 1 1
7 0 1 1 1 1 0 1 1 0 1 0 0 0 1
8 1 0 0 0 1 0 1 0 0 0 1 1 0 0 1 0
2 1 1 1 1
3 1 1 1 0 1 0
9 1 0 0 0 1 1 1 1 1 1 0 0 1 0 0 1 0 0
2 0 1 0 0
4 1 0 1 1 1 1 0 1
10 1 1 1 1 0 1 1 0 0 0 0 1 0 1 0 0 0 1 0 1
9 0 1 0 1 0 0 1 1 1 1 1 0 0 0 1 0 0 1
4 0 1 1 1 1 1 1 1
4 0 1 0 0 1 1 1 0
7 1 1 1 0 0 0 0 0 0 1 0 1 0 1
9 1 0 0 0 1 0 0 1 1 1 1 0 1 0 1 1 1 1
5 0 0 0 0 0 1 0 0 1 0
9 1 0 0 1 1 0 1 0 1 0 0 0 1 1 0 1 1 1
10 0 0 0 1 1 0 0 1 0 0 0 0 0 1 0 0 0 0 1 1
10 0 1 0 0 1 1 1 1 0 0 0 1 1 1 0 1 0 1 0 1
5 1 0 0 0 1 0 1 1 0 1
6 0 1 0 1 0 0 1 0 1 0 0 0
1 0 0
6 0 0 1 1 0 0 1 1 0 1 1 0
6 1 1 0 1 1 1 0 0 1 1 0 0
9 0 0 0 1 0 0 0 0 0 0 0 1 0 1 0 0 1 1
7 0 0 1 0 1 1 0 1 0 1 0 1 1 1
8 1 0 1 1 1 0 1 0 1 1 0 1 1 1 0 0
5 0 1 1 0 1 1 1 0 0 1
4 1 0 0 1 0 0 0 1
6 1 1 1 1 1 1 0 0 0 0 1 0
6 0 1 1 0 0 1 1 0 0 0 0 0
9 1 1 1 1 1 0 1 1 1 1 1 0 1 0 0 0 1 1
8 0 1 0 0 1 0 0 1 1 0 0 0 1 0 1 0
10 1 0 1 0 0 0 1 0 1 1 1 1 1 1 1 0 1 0 1 0
5 0 1 1 1 1 0 1 1 0 1
6 0 1 1 0 0 0 1 1 0 0 0 0
`

var expectedOutputs = []string{
	"3",
	"1",
	"0",
	"1",
	"1",
	"4",
	"0",
	"4",
	"0",
	"2",
	"0",
	"1",
	"2",
	"0",
	"3",
	"3",
	"1",
	"0",
	"2",
	"0",
	"0",
	"1",
	"2",
	"2",
	"1",
	"1",
	"1",
	"1",
	"1",
	"1",
	"1",
	"1",
	"0",
	"1",
	"1",
	"3",
	"2",
	"0",
	"1",
	"2",
	"0",
	"1",
	"1",
	"1",
	"1",
	"1",
	"1",
	"3",
	"0",
	"1",
	"1",
	"2",
	"2",
	"1",
	"1",
	"2",
	"2",
	"1",
	"2",
	"3",
	"1",
	"1",
	"1",
	"2",
	"2",
	"1",
	"3",
	"1",
	"0",
	"2",
	"4",
	"1",
	"1",
	"3",
	"2",
	"1",
	"2",
	"1",
	"4",
	"2",
	"1",
	"1",
	"2",
	"2",
	"1",
	"0",
	"3",
	"4",
	"4",
	"2",
	"1",
	"1",
	"1",
	"5",
	"3",
	"3",
	"1",
	"3",
	"1",
	"1",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: empty", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", i+1, 2*n, len(fields)-1)
		}
		a := make([]int, n)
		b := make([]int, n)
		for idx := 0; idx < n; idx++ {
			val, err := strconv.Atoi(fields[1+idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad a value: %v", i+1, err)
			}
			a[idx] = val
		}
		for idx := 0; idx < n; idx++ {
			val, err := strconv.Atoi(fields[1+n+idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad b value: %v", i+1, err)
			}
			b[idx] = val
		}
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
