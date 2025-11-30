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
	arr []int64
}

// solve embeds the logic from 1730E.go.
func solve(tc testCase) string {
	n := tc.n
	arr := tc.arr
	var ans int64
	for l := 0; l < n; l++ {
		minv := arr[l]
		maxv := arr[l]
		for r := l; r < n; r++ {
			if arr[r] < minv {
				minv = arr[r]
			}
			if arr[r] > maxv {
				maxv = arr[r]
			}
			if maxv%minv == 0 {
				ans++
			}
		}
	}
	return strconv.FormatInt(ans, 10)
}

// Embedded copy of testcasesE.txt.
const testcaseData = `
7 3 1 5 8 10 2 6
2 2 8
4 9 1 1 10
2 8 5
5 9 7 3 4 2
3 8 4 1
4 3 10 3 3
3 2 2 5
2 5 10
8 10 5 10 7 1 2 3 8
1 6
6 5 7 7 1 5 4
4 2 10 3 4
4 8 6 1 5
7 10 9 9 3 6 10 9
7 4 5 2 8 5 4 5
5 2 6 3 8 1
7 8 1 5 9 10 7 8
4 3 1 5 2
7 5 4 8 8 5 9 7
3 10 8 5
3 1 4 6
5 2 1 7 5 6
4 6 7 3 2
4 6 8 8 3
6 7 7 7 6 10 2
2 9 8
3 2 1 9
7 8 4 7 2 7 1 5
3 5 5 9
2 9 6
1 3
1 2
6 6 8 4 4 6 10
8 5 9 3 8 7 1 2 1
2 6 4
5 9 2 9 5 2
2 9 10
5 1 4 2 3 8
3 3 8 5
3 8 10 1
7 1 2 2 5 2 7 4
5 7 4 8 10 9
6 3 7 1 10 7 10
4 5 6 9 6
8 8 5 6 5 3 3 4 8
6 1 6 8 10 5 4
1 2
8 9 3 4 7 10 6 9 3
1 8
7 8 8 6 1 4 10 4
2 8 3
2 9 9
7 10 5 4 2 7 1 9
8 6 1 4 10 5 7 8 9
5 2 1 2 6 7
7 7 4 5 8 7 9 1
7 5 8 3 7 7 7 10
5 8 9 6 5 7
3 6 1 7
8 8 5 2 4 6 4 7 1
6 3 3 8 2 2 3
1 4
3 1 9 7
7 3 5 7 8 5 9 1
3 8 3 2
2 8 5
5 5 7 1 10 5
6 8 2 1 5 4 8
5 9 4 7 2 7
1 1
2 6 1
5 5 1 1 1 10
1 3
3 8 2 10
3 1 9 3
2 2 2
2 4 2
8 7 2 2 8 8 9 3 9
2 3 6
7 9 3 2 9 1 9 4
2 7 6
6 9 10 9 7 5 4
2 3 10
7 6 5 6 9 6 1 1
1 3
5 4 3 3 2 8
8 3 5 5 3 1 3 6 6
4 2 9 8 3
5 6 4 2 10 4
6 7 2 5 7 9 6
1 6
4 7 1 7 9
5 3 5 1 5 2
6 7 1 5 4 3 5
6 9 2 3 8 4 10
2 7 5
1 5
1 1
3 6 3 7
`

var expectedOutputs = []string{
	"26",
	"3",
	"10",
	"2",
	"9",
	"6",
	"5",
	"4",
	"3",
	"33",
	"1",
	"18",
	"7",
	"9",
	"13",
	"21",
	"13",
	"21",
	"9",
	"14",
	"4",
	"5",
	"12",
	"4",
	"5",
	"14",
	"2",
	"6",
	"22",
	"4",
	"2",
	"1",
	"1",
	"13",
	"32",
	"2",
	"5",
	"2",
	"13",
	"3",
	"5",
	"14",
	"7",
	"17",
	"4",
	"16",
	"14",
	"1",
	"13",
	"1",
	"23",
	"2",
	"3",
	"22",
	"25",
	"13",
	"17",
	"10",
	"5",
	"6",
	"25",
	"17",
	"1",
	"5",
	"14",
	"4",
	"2",
	"14",
	"20",
	"5",
	"1",
	"3",
	"15",
	"1",
	"6",
	"6",
	"3",
	"3",
	"23",
	"3",
	"22",
	"2",
	"8",
	"2",
	"18",
	"1",
	"11",
	"31",
	"5",
	"13",
	"6",
	"1",
	"9",
	"13",
	"15",
	"11",
	"2",
	"1",
	"1",
	"4",
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
		if len(parts) < 1 {
			return nil, fmt.Errorf("line %d: empty", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(parts) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", i+1, n, len(parts)-1)
		}
		arr := make([]int64, n)
		for idx := 0; idx < n; idx++ {
			v, err := strconv.ParseInt(parts[1+idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d: %v", i+1, idx, err)
			}
			arr[idx] = v
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(v, 10))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
