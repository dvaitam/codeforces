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
	c int
	a []int
}

// solve embeds the logic from 1730A.go.
func solve(tc testCase) string {
	counts := make(map[int]int)
	for _, x := range tc.a {
		counts[x]++
	}
	ans := 0
	for _, cnt := range counts {
		if cnt < tc.c {
			ans += cnt
		} else {
			ans += tc.c
		}
	}
	return strconv.Itoa(ans)
}

// Embedded copy of testcasesA.txt.
const testcaseData = `
7 4 1 5 9 8 7 5 8
6 5 4 9 3 5 3 2
10 3 9 10 3 5 2 2 6 8 9 2
6 4 6 10 4 9 8 8
9 3 1 9 1 2 7 1 10 8 6
4 3 2 4 10 4
4 2 9 8 2 2
6 5 8 2 5 9 5 2
9 3 9 4 10 9 10 5 8 2 10
7 3 10 4 5 3 4 3 1
10 3 8 2 2 3 3 1 2 9 7 9
5 5 4 4 10 7 10
5 4 8 6 2 6 10
2 4 10 6
4 2 1 5 2 4
6 2 6 7 1 2 3 4
1 5 9
10 1 1 2 4 10 10 2 7 2 6 2
1 5 1
4 2 2 8 4 1
1 5 7
10 1 5 2 4 2 5 6 7 3 1 9
8 1 10 2 7 4 5 6 8 10
3 2 1 3 3
6 5 5 2 10 8 3 1
8 4 10 9 5 6 7 5 3 9
1 4 2
6 1 9 5 3 4 8 6
10 3 6 10 10 3 5 7 7 2 1 10
4 3 3 4 4 8
7 5 7 1 7 10 7 1 3
8 1 5 3 8 9 8 9 10 1
1 4 6
5 4 1 7 4 9 2
3 1 7 7 6
1 2 1
1 5 10
2 2 2 10
4 3 5 3 2 8
7 1 1 5 8 2 5 3 9
6 1 3 5 1 1 1 4
5 5 6 6 10 1 10
8 4 7 6 9 3 4 7 10 5
1 2 3
5 3 6 6 2 6 10
1 1 5
3 2 10 5 6
7 5 3 5 2 8 4 1 5
3 5 2 5 7
6 3 7 2 2 9 8 8
6 3 2 8 2 8 7 1
5 3 3 3 10 7 2
2 1 4 4
1 4 1
2 4 9 9
5 4 8 10 4 7 2
6 2 5 10 3 7 4 6
2 1 1 9
8 2 2 8 7 5 4 1 4 10
3 1 4 8 7
6 5 3 2 10 8 3 10
7 4 9 8 6 8 8 4 9
10 2 1 6 6 6 1 9 3 5 10 3
7 5 5 8 2 2 9 1 2
4 2 1 5 1 8
6 2 3 8 6 9 7 9
9 1 10 2 9 10 2 7 4 5 9
10 4 8 7 10 10 4 1 1 3 5 9
10 3 6 2 8 5 5 7 7 7 1 3
3 2 5 6 1
1 4 7
3 4 10 2 3
6 4 1 10 8 7 8 1
2 4 3 1
1 5 10
3 3 2 9 6
4 4 8 2 1 10
8 5 6 2 10 5 3 7 5 2
9 2 1 7 8 6 4 8 6 2 1
1 4 5
1 5 10
10 2 4 2 9 9 7 9 5 2 3 7
10 4 2 2 7 2 2 7 3 1 8 7
7 1 8 6 5 2 6 2 2
6 1 6 6 3 1 4 6
2 5 3 4
1 2 2
1 3 6
1 5 4
3 2 8 2 8
6 3 3 1 4 6 6 8
5 3 9 6 3 10 2
2 5 10 5
3 4 3 3 4
6 5 4 4 3 5 6 7
1 2 10
1 4 2
2 2 7 5
9 4 3 10 7 5 6 2 4 8 6
9 1 7 7 1 7 6 8 4 6 5
`

var expectedOutputs = []string{
	"7",
	"6",
	"10",
	"6",
	"9",
	"4",
	"4",
	"6",
	"9",
	"7",
	"10",
	"5",
	"5",
	"2",
	"4",
	"6",
	"1",
	"6",
	"1",
	"4",
	"1",
	"8",
	"7",
	"3",
	"6",
	"8",
	"1",
	"6",
	"10",
	"4",
	"7",
	"6",
	"1",
	"5",
	"2",
	"1",
	"1",
	"2",
	"4",
	"6",
	"4",
	"5",
	"8",
	"1",
	"5",
	"1",
	"3",
	"7",
	"3",
	"6",
	"6",
	"5",
	"1",
	"1",
	"2",
	"5",
	"6",
	"2",
	"8",
	"3",
	"6",
	"7",
	"9",
	"7",
	"4",
	"6",
	"6",
	"10",
	"10",
	"3",
	"1",
	"3",
	"6",
	"2",
	"1",
	"3",
	"4",
	"8",
	"9",
	"1",
	"1",
	"9",
	"10",
	"4",
	"4",
	"2",
	"1",
	"1",
	"1",
	"3",
	"6",
	"5",
	"2",
	"3",
	"6",
	"1",
	"1",
	"2",
	"9",
	"6",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		c, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad c: %v", i+1, err)
		}
		if len(parts) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers for a, got %d", i+1, n, len(parts)-2)
		}
		arr := make([]int, n)
		for idx := 0; idx < n; idx++ {
			v, err := strconv.Atoi(parts[2+idx])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d: %v", i+1, idx, err)
			}
			arr[idx] = v
		}
		tests = append(tests, testCase{n: n, c: c, a: arr})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.c))
	for i, v := range tc.a {
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
