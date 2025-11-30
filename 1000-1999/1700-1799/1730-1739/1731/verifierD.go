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
	m int
	a [][]int
}

// solve embeds the logic from 1731D.go.
func solve(tc testCase) string {
	matrix := tc.a
	n, m := tc.n, tc.m
	check := func(l int) bool {
		pref := make([][]int, n+1)
		for i := range pref {
			pref[i] = make([]int, m+1)
		}
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				val := 0
				if matrix[i-1][j-1] >= l {
					val = 1
				}
				pref[i][j] = pref[i-1][j] + pref[i][j-1] - pref[i-1][j-1] + val
			}
		}
		target := l * l
		for i := l; i <= n; i++ {
			for j := l; j <= m; j++ {
				sum := pref[i][j] - pref[i-l][j] - pref[i][j-l] + pref[i-l][j-l]
				if sum == target {
					return true
				}
			}
		}
		return false
	}

	low, high := 1, n
	if m < high {
		high = m
	}
	for low < high {
		mid := (low + high + 1) / 2
		if check(mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return strconv.Itoa(low)
}

// Embedded copy of testcasesD.txt.
const testcaseData = `
4 4 2 4 5 5 3 7 4 4 10 3 5 7 10 1 9 8
1 5 9 5 2 10 7
2 4 8 6 4 10 5 3 8 3
5 6 7 3 1 9 4 5 5 7 1 10 8 7 6 1 2 5 7 4 3 6 1 7 8 5 9 10 5 6 9 2
2 4 8 10 9 10 10 8 10 2
3 5 7 4 7 2 4 1 10 4 10 1 10 3 9 4 6
5 5 3 2 8 5 2 4 5 5 4 6 1 1 5 8 7 7 5 1 6 6 5 6 10 4 6
4 6 9 8 1 10 10 8 1 10 6 2 1 5 10 3 6 9 3 4 6 6 2 7 4 4
5 5 5 10 6 3 8 5 10 2 1 2 9 6 3 5 3 4 3 6 3 6 6 9 4 3 1
3 4 7 8 9 5 8 9 8 4 4 8 5 3
2 3 1 2 8 6 3 2
1 6 7 9 7 5 3 4
4 5 3 4 2 8 6 6 9 6 4 8 4 6 9 8 3 4 2 6 1 7
5 6 3 7 3 9 2 6 3 3 4 10 10 10 8 8 10 7 3 4 10 5 9 3 9 7 9 8 1 6 6 2
5 5 2 3 4 9 9 6 9 3 5 8 2 3 7 3 6 10 7 6 8 2 6 9 2 8 6
3 5 2 9 4 2 1 5 5 6 4 3 6 2 9 7 2
4 4 10 1 6 9 4 2 2 7 8 1 4 8 5 10 9 1
2 3 2 6 10 4 7 1
3 5 8 9 5 7 9 5 3 5 10 5 5 10 2 5 10
4 6 5 2 5 4 6 7 7 1 2 7 5 10 1 3 6 8 8 3 8 1 3 9 10 10
5 6 2 1 3 3 4 10 3 9 4 8 8 1 4 3 4 4 1 7 2 9 4 1 6 5 2 7 1 7 3 8
1 5 10 8 10 10 5
5 6 4 2 3 6 3 3 5 10 8 3 4 6 1 9 9 7 10 4 6 7 10 3 2 1 8 9 2 6 4 5
1 4 3 7 5 7
3 3 5 10 2 10 2 6 2 9 1
1 6 2 6 7 9 10 2
5 5 8 6 6 2 3 8 4 10 9 9 8 5 8 5 7 8 8 6 5 5 1 7 4 9 1
5 5 5 10 7 1 5 9 7 6 5 9 7 7 7 9 10 3 10 6 10 9 2 10 4 7 4
3 5 7 8 4 4 8 1 5 1 8 10 5 7 1 2 8
2 5 5 2 9 9 2 8 10 1 10 5
3 4 6 5 5 9 2 2 9 3 10 2 10 10
4 6 1 1 8 6 3 8 6 7 2 5 10 10 2 2 1 1 7 2 2 3 10 4 5 9
1 4 5 7 6 5
1 1 6
1 4 10 4 2 10
5 6 3 10 7 3 6 6 7 3 9 5 8 4 1 9 4 5 6 6 10 1 9 2 5 3 10 5 1 10 10 3
5 6 2 1 3 3 8 1 3 6 5 1 1 1 10 1 10 4 5 3 4 4 8 3 4 2 7 3 8 6 3 7
5 5 8 2 2 10 4 1 9 9 1 1 9 7 7 7 8 6 8 4 7 1 6 10 8 3 1
1 6 9 2 5 5 7 1
5 6 5 4 4 7 2 10 5 9 3 4 6 9 7 3 9 8 9 5 1 4 2 10 5 10 6 3 8 10 7 8
4 4 5 4 6 7 4 6 1 7 7 3 8 6 3 1 2 6
2 2 9 9 2 5
5 5 3 7 9 2 1 9 3 6 1 3 8 7 2 6 10 1 5 5 1 8 10 4 8 5 9
3 3 1 5 2 5 10 10 9 6 8
1 6 10 3 9 1 4 3
1 3 6 3 8
1 1 10
3 3 7 2 6 2 4 3 10 5 9
1 6 3 10 4 3 6 3
4 5 7 9 5 6 9 1 8 10 1 1 3 9 1 5 2 6 4 6 9 10
3 5 9 6 4 10 4 7 1 5 6 4 7 4 3 2 3
1 3 8 5 2
1 1 2
4 6 8 1 3 1 7 3 1 10 6 3 5 9 3 1 6 5 7 5 2 5 9 4 4 1
3 4 5 7 4 5 5 9 6 4 5 3 7 7
3 3 7 6 3 1 8 5 3 9 10
5 5 1 9 8 1 3 1 5 9 7 3 2 9 10 9 6 9 2 7 5 1 6 9 3 6 2
1 2 1 6
2 3 7 2 9 9 1 1
5 5 9 1 4 10 5 7 7 8 1 3 9 1 5 2 4 9 7 9 10 3 5 2 6 10 5
5 5 2 1 6 8 9 5 8 10 3 3 6 2 7 10 7 6 1 7 5 10 7 4 4 8 10
4 5 8 10 8 7 7 1 2 5 4 6 6 6 9 1 7 7 9 4 1 5
3 3 4 10 7 4 9 9 5 7 7
5 5 9 1 1 8 6 8 7 6 6 8 6 10 6 9 5 4 9 5 8 10 8 7 8 2 10
5 6 2 10 10 10 3 4 8 2 2 1 3 6 5 3 5 5 7 6 10 9 10 8 3 3 10 6 2 6 9 5
1 3 5 2 5
3 4 6 5 3 8 4 8 1 3 4 4 5 4
3 3 3 9 8 7 1 2 10 7 9
4 5 10 10 4 2 6 1 7 5 4 5 7 3 8 6 1 9 9 5 7 6
1 4 6 1 1 3
2 4 9 9 6 7 7 4 9 3
2 4 3 10 5 4 7 5 9 2
5 6 8 6 3 4 7 2 10 5 2 6 5 7 7 7 5 4 7 5 6 3 10 2 8 8 1 8 3 6 3 2
5 5 10 6 2 10 9 2 5 9 7 1 3 10 9 4 4 2 10 7 4 2 2 9 1 6 1
3 4 10 5 9 1 2 9 9 5 6 2 6 10
5 5 3 5 7 5 1 3 4 10 3 1 1 5 9 1 3 1 9 6 9 10 9 7 7 4 4
1 4 4 6 9 3
3 6 3 10 6 6 10 6 8 2 10 5 3 5 5 9 3 1 10 5
1 4 3 10 9 10
5 5 10 8 2 8 8 10 4 1 1 8 5 5 3 3 7 9 5 9 5 5 5 8 7 8 2
5 6 3 6 4 7 3 5 10 10 5 7 6 8 10 3 8 8 10 6 10 5 3 4 2 6 8 3 7 4 3 8
2 6 10 5 9 1 4 5 10 10 3 9 6 4
2 5 9 4 9 10 1 7 6 2 9 8
1 2 4 8
1 5 8 5 10 10 4
5 6 7 7 6 5 8 7 5 4 8 10 7 2 10 4 1 7 10 7 9 8 6 5 4 1 9 2 2 9 5 8
4 6 1 10 2 3 4 2 7 6 9 3 10 4 7 10 10 9 2 1 4 5 10 6 6 4
4 6 9 9 8 7 8 2 10 10 9 10 4 1 1 9 7 8 5 5 8 6 7 5 9 5
1 3 4 2 7
5 5 2 2 10 2 5 10 2 7 10 10 10 9 3 1 5 10 3 7 5 2 1 4 5 4 9
4 4 3 3 7 1 9 3 4 4 1 4 4 2 3 2 3 9
1 4 7 5 4 10
5 5 5 3 1 9 8 6 8 8 4 3 1 9 6 1 1 9 5 2 1 5 6 4 4 4 3
2 2 7 8 6 7
1 2 2 2
4 6 8 2 10 9 8 8 5 9 1 6 5 3 2 10 9 5 10 5 5 8 7 6 9 10
1 2 4 10
4 6 3 6 9 1 2 10 7 8 8 8 3 3 6 10 2 1 5 9 4 7 8 8 7 3
3 4 4 8 9 6 7 10 3 3 1 2 2 5
2 6 4 2 10 1 7 9 1 7 1 8 10 5`

var expectedOutputs = []string{
	"3",
	"1",
	"2",
	"3",
	"2",
	"2",
	"2",
	"2",
	"3",
	"3",
	"2",
	"1",
	"2",
	"3",
	"3",
	"2",
	"2",
	"2",
	"2",
	"3",
	"2",
	"1",
	"3",
	"1",
	"2",
	"1",
	"3",
	"4",
	"2",
	"2",
	"2",
	"2",
	"1",
	"1",
	"1",
	"3",
	"3",
	"3",
	"1",
	"3",
	"2",
	"2",
	"2",
	"2",
	"1",
	"1",
	"1",
	"2",
	"1",
	"2",
	"2",
	"1",
	"1",
	"3",
	"3",
	"2",
	"2",
	"1",
	"1",
	"2",
	"3",
	"2",
	"3",
	"3",
	"3",
	"1",
	"2",
	"1",
	"3",
	"1",
	"2",
	"2",
	"2",
	"3",
	"2",
	"2",
	"1",
	"2",
	"1",
	"3",
	"3",
	"2",
	"2",
	"1",
	"1",
	"2",
	"3",
	"4",
	"1",
	"2",
	"2",
	"1",
	"2",
	"2",
	"1",
	"3",
	"1",
	"2",
	"2",
	"2",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		vals := strings.Fields(line)
		if len(vals) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		vals = vals[2:]
		if len(vals) != n*m {
			return nil, fmt.Errorf("line %d: expected %d grid values, got %d", i+1, n*m, len(vals))
		}
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			row := make([]int, m)
			for c := 0; c < m; c++ {
				v, err := strconv.Atoi(vals[r*m+c])
				if err != nil {
					return nil, fmt.Errorf("line %d: bad value: %v", i+1, err)
				}
				row[c] = v
			}
			grid[r] = row
		}
		tests = append(tests, testCase{n: n, m: m, a: grid})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.a {
		for i, v := range row {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
	}

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
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
