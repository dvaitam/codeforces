package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n, m int
	grid [][]int
}

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `2 2 0 1 0 1
4 6 0 0 0 1 1 1 1 1 0 0 1 1 1 1 1 0 0 0 0 0 0 1 0 0
6 6 1 0 1 1 1 1 1 1 0 1 1 0 1 1 1 1 1 1 1 1 1 0 1 0 1 1 1 1 1 1 0 1 1 0 1 0
3 2 0 0 1 0 0 0
4 3 0 0 1 0 0 1 1 0 0 0 0 0
2 2 0 0 1 1
3 3 0 0 1 0 0 0 0 0 1
6 2 1 1 1 0 1 1 0 1 1 0 1 0
2 4 0 0 1 0 1 1 1 0
4 4 1 1 0 0 0 1 0 0 0 0 0 1 0 0 0 0
5 2 1 0 0 1 1 1 1 0 0 0
5 5 0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 1 0 0 1 1 0 0 1 0 0
5 4 0 0 1 1 1 1 0 1 0 0 1 0 0 1 0 1 1 1 0 0
2 4 0 1 1 0 0 1 0 0
6 4 1 1 0 1 1 0 1 0 0 0 1 1 0 0 0 1 0 1 1 1 1 0 1 1
5 6 1 1 1 1 1 0 0 1 1 1 0 0 0 1 0 0 1 0 0 0 0 1 1 0 1 1 0 1 0 0
6 5 0 1 0 1 0 0 1 0 1 1 0 1 1 0 1 0 1 0 1 1 1 1 1 1 1 1 0 0 0 1
5 2 1 1 1 0 0 0 1 0 1 0
2 5 0 1 0 1 0 0 1 1 1 0
5 3 0 1 0 1 0 1 1 1 1 1 0 0 1 0 1
4 3 0 0 1 1 1 1 0 0 1 1 0 1
2 6 1 1 0 1 1 0 1 1 0 1 1 1
4 4 1 0 0 1 0 0 1 0 0 1 0 1 0 0 0 1
4 3 1 1 0 1 1 1 0 1 1 1 1 0
5 3 1 1 0 0 0 1 0 1 0 1 0 0 0 1 1
5 4 1 1 1 1 0 1 1 1 1 0 1 1 1 1 0 0 1 1 0 0
2 4 0 1 0 1 0 1 0 1
6 4 1 0 1 1 0 1 0 0 0 1 1 0 1 1 1 0 1 0 0 0 0 1 1 0
3 3 0 1 1 1 0 1 0 1 1
3 2 1 1 1 1 0 1
5 5 0 1 1 1 0 0 1 0 0 0 0 1 1 1 0 1 0 1 0 0 0 0 1 0 1
5 2 1 0 0 1 0 0 0 0 0 0
6 2 1 0 0 1 1 0 0 0 0 0 1 0
2 6 0 1 1 1 0 0 0 0 0 1 0 1
4 5 0 1 0 1 1 1 0 0 1 1 1 1 1 0 1 0 0 0 0 0
3 5 0 0 0 1 1 0 0 0 1 0 0 0 0 0 0
5 4 0 0 1 1 0 1 1 0 1 1 0 1 0 0 1 1 1 1 0 1
3 3 1 1 1 1 1 1 0 1 1
3 5 1 1 0 0 1 1 0 0 1 0 1 0 0 0 0
4 2 1 1 1 1 0 1 1 0
4 6 0 0 0 0 1 0 1 1 0 1 1 0 0 1 0 0 1 1 1 0 0 0 1 0
3 6 1 1 0 1 1 1 1 1 1 0 0 0 0 0 0 0 1 1
5 3 1 0 1 1 1 0 1 0 1 0 0 1 0 0 0
3 4 0 0 0 1 0 0 0 1 0 0 1 1
5 2 1 1 0 1 0 1 1 1 1 1
6 6 0 0 0 1 1 0 1 0 0 1 1 1 0 0 0 0 0 1 0 1 0 0 0 0 1 0 1 0 0 0 0 0 0 0 1 0
6 4 0 0 1 1 0 0 1 0 1 1 0 0 1 1 1 0 1 0 1 0 0 1 0 0
5 6 1 1 1 1 0 1 1 1 0 1 0 1 0 1 1 0 1 1 1 0 1 1 1 1 0 0 1 1 1 1
2 6 0 1 0 1 1 1 0 1 1 1 1 0
3 3 0 1 1 1 0 1 1 1 1
4 4 1 0 1 0 0 1 0 0 0 1 1 1 1 0 1 0
6 6 1 1 0 1 0 0 1 0 1 1 0 0 0 1 1 1 0 1 0 0 1 0 1 1 1 1 1 1 0 0 1 0 1 1 0 0
6 6 0 0 1 0 1 1 0 0 0 1 1 1 0 0 0 1 0 0 0 0 0 0 1 0 0 1 1 0 0 1 0 1 0 1 1 1
2 2 1 0 0 1
2 4 0 1 0 1 1 0 0 0
5 3 0 1 1 0 0 1 1 0 0 0 1 1 0 0 0
3 2 1 0 1 0 1 0
6 6 0 1 0 1 0 1 0 0 1 0 0 0 0 1 1 1 1 1 1 0 0 0 0 0 0 1 0 1 0 0 0 0 1 0 0 1
3 5 0 1 1 1 1 0 0 0 0 1 0 1 0 1 0
2 6 0 0 0 1 0 0 0 1 1 1 0 1
3 2 1 0 1 0 0 1
4 6 0 1 1 0 0 0 0 0 0 1 0 0 1 1 1 0 0 1 0 0 0 0 1 0
3 5 1 0 0 1 1 1 0 1 1 1 1 1 1 1 1
3 5 1 1 1 1 0 0 1 1 1 0 1 1 0 0 1
2 3 0 1 0 0 1 0
2 4 1 0 0 1 0 0 0 1
6 2 1 1 0 0 1 0 0 0 1 1 1 0
2 6 1 0 0 1 0 1 1 0 1 1 0 1
2 3 1 1 1 1 1 0
4 2 1 0 0 1 1 0 1 1
2 3 0 1 0 0 0 0
4 3 1 1 1 0 0 0 0 0 0 0 0 1
4 2 0 0 1 1 0 1 0 1
2 3 0 0 1 0 1 0
3 6 0 0 0 0 0 0 0 0 0 0 1 1 0 0 0 0 0 1
5 3 0 0 1 1 0 0 0 1 1 1 1 0 1 0 1
6 3 1 1 1 1 0 0 0 0 0 1 0 0 1 0 0 1 0 1
4 6 0 1 0 0 1 1 1 1 0 0 1 1 1 0 0 0 0 0 1 0 1 0 0 0
2 2 0 0 1 0
5 6 1 0 0 0 0 0 0 0 1 1 0 0 1 0 0 0 0 0 0 1 1 0 0 0 1 0 0 1 0 0
3 2 1 1 1 0 1 0
4 5 0 0 0 1 1 0 1 1 1 0 1 0 1 1 0 1 0 1 1 0
5 4 1 1 0 0 0 0 0 1 0 1 1 0 0 0 0 0 0 1 0 0
5 4 0 1 0 0 1 1 0 0 0 1 0 0 0 1 1 1 0 0 1 1
6 6 0 1 0 1 1 0 0 1 0 0 0 0 1 0 1 1 1 0 0 0 0 1 1 0 0 0 0 1 0 0 1 1 0 0 1 0
4 3 0 0 1 1 1 1 1 1 0 1 0 0
6 6 0 1 1 1 1 1 1 0 1 0 0 0 1 0 0 0 1 1 1 0 0 1 1 0 0 1 1 1 1 1 1 1 1 1 0 1
2 6 0 0 1 1 1 0 1 1 1 1 1 0
6 5 1 1 0 0 0 0 1 0 1 0 1 1 1 0 1 1 1 1 1 0 1 0 1 0 0 0 0 0 1 1
6 3 1 1 0 0 0 0 0 1 0 1 1 0 1 0 1 0 1 1
4 2 0 1 0 1 1 1 1 0
5 5 0 1 1 0 1 1 0 1 1 1 1 0 1 0 0 0 0 0 1 1 0 0 0 0 1
6 5 1 0 0 0 0 1 1 1 1 1 1 0 1 1 0 1 1 0 0 1 1 1 1 0 1 1 1 0 0 1
6 6 0 1 1 0 1 0 1 0 0 0 1 0 0 0 0 1 0 1 0 1 1 1 0 0 1 0 1 0 0 0 0 1 0 1 0 0
4 2 0 1 1 0 0 0 1 0
5 6 0 1 0 0 0 0 1 0 0 1 1 1 1 0 1 1 0 0 0 1 0 0 0 0 1 1 1 0 1 0
5 6 1 0 1 1 0 0 1 1 0 1 1 1 0 0 1 1 1 0 0 0 1 1 1 1 1 1 1 1 0 0
3 3 0 0 0 0 1 0 1 1 0
4 5 0 0 0 1 1 0 0 1 0 0 0 0 0 0 0 0 1 0 0 1
2 6 0 1 0 1 0 0 0 0 1 0 1 1`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		fields := strings.Fields(ln)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: expected at least 2 fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		expected := 2 + n*m
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, expected, len(fields))
		}
		grid := make([][]int, n)
		pos := 2
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				v, err := strconv.Atoi(fields[pos])
				if err != nil {
					return nil, fmt.Errorf("line %d: parse cell %d: %v", idx+1, pos-1, err)
				}
				row[j] = v
				pos++
			}
			grid[i] = row
		}
		cases = append(cases, testcase{n: n, m: m, grid: grid})
	}
	return cases, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// solve implements the logic from 1366C.go for a single testcase.
func solve(tc testcase) string {
	n, m := tc.n, tc.m
	zeros := make([]int, n+m-1)
	ones := make([]int, n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if tc.grid[i][j] == 0 {
				zeros[i+j]++
			} else {
				ones[i+j]++
			}
		}
	}
	total := n + m - 2
	ans := 0
	for l, r := 0, total; l < r; l, r = l+1, r-1 {
		ans += min(zeros[l]+zeros[r], ones[l]+ones[r])
	}
	return strconv.Itoa(ans)
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(strconv.Itoa(tc.grid[i][j]))
			}
			input.WriteByte('\n')
		}
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solve(tc))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
