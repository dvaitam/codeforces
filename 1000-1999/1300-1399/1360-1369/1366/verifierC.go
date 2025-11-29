package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
2 5 1 0 0 1 1 0 0 0 0 0
6 2 0 0 0 0 0 1 1 0 1 0 0 0
2 2 0 1 1 1
6 6 0 0 0 0 1 1 1 0 1 0 1 0 0 1 1 0 0 1 0 1 1 1 1 0 1 1 0 1 1 0 0 1 0 1 0 0
5 1 1 0 0 1 0 0
4 2 1 1 0 1 1 0 0 0 0 0
6 6 0 0 0 0 1 1 1 1 1 1 0 0 0 0 1 1 1 1 0 0 0 0 0 1 1 0 1 0 1 1 0 0 1 1 1 1
4 2 0 0 1 1 0 1 1 0 0 1
3 1 0 0 0 0
6 6 0 1 0 0 1 1 1 1 1 0 1 0 0 0 0 0 0 0 0 0 1 1 1 1 1 0 1 0 0 1 1 0 0 1 0 0
3 5 1 1 0 0 0 0 1 1 0 1 0 0 0 0 0 0
6 4 0 0 1 1 0 0 1 1 0 0 0 0 1 0 1 1 0 1 0 0 1 1 1 1
6 3 1 0 1 0 0 0 1 1 0 1 1 0 0 1 0 1 0
6 6 1 1 0 0 1 1 0 0 1 0 1 0 0 1 0 1 1 1 0 1 1 1 0 1 1 1 0 0 1 0 1 0 0 0
3 2 0 1 0 1 1 1
6 2 1 1 0 0 1 0 0 1 1 0 1 0
2 6 1 0 1 1 1 1 1 1 0 0 0 1 0 0
4 5 1 0 1 0 1 0 1 1 0 0 1 1 0 0 1 0 1 0 0 0
6 1 0 0 1 0 0 1 1
6 3 1 1 1 1 1 0 0 0 1 1 0 1 1 1 0 0 0
2 1 1 1 0
2 5 1 1 1 0 1 0 1 1 1 1
5 3 0 1 1 0 0 1 1 0 1 0 1 1 0 0 0 1 1
2 4 1 0 1 1 0 0 1 1
4 6 0 0 1 1 0 1 1 1 1 0 1 0 0 1 1 1 0 0 0 1 0 1 0 0 1 1
2 6 1 0 1 0 1 0 0 0 0 1 0 1 1 0
5 5 0 0 1 1 1 1 1 1 0 1 0 0 1 0 1 0 0 1 1 1 0 0 1 0 1
5 3 1 1 0 0 1 1 0 1 1 1 0 0 0 0 1
5 2 0 1 0 0 0 1 1 0 1 0
5 2 0 0 0 0 0 1 0 1 1 0
6 5 1 0 1 0 0 1 1 0 1 0 0 1 1 1 0 0 0 1 0 0 0 1 0 0 1 0 1 0 0 1
5 5 1 1 0 0 0 1 1 1 0 0 1 1 0 0 0 1 0 1 0 1 0 1 1 0 1
2 6 1 0 0 0 1 1 0 0 1 1 1 1 1 0
3 6 0 0 1 1 0 1 1 0 0 0 1 1 0 0 1 1 0 1 0
5 1 0 0 1 0 0
6 3 1 1 0 1 1 1 1 1 0 1 1 1 0 1 1 1 0
5 4 1 0 1 0 0 1 1 1 0 1 1 0 1 0 1 0 1 0 1
5 1 1 1 0 0 0
2 6 1 1 1 0 0 0 0 1 1 1 1 1 0
5 2 1 0 1 1 0 1 0 1 0 1 0
2 6 1 0 0 1 1 0 1 0 0 1 0 1 0 0
5 6 1 1 0 1 0 1 1 1 0 1 1 1 1 1 0 1 1 1 1 0 0 1 0 1 1 1 1 0 1 1 1 0 1
2 6 0 1 1 1 0 0 0 1 0 1 0 0 1 0
6 4 1 0 0 0 0 1 0 0 0 0 1 0 0 1 1 1 1 0 1 0 1 1 1 1
5 5 1 1 0 1 0 0 1 1 1 0 1 1 1 1 0 0 0 0 0 0 0 0 1 0 0
2 2 1 0 0 1
6 6 0 0 1 0 0 0 1 1 1 0 0 1 1 1 0 0 1 0 1 1 0 1 0 1 1 1 0 1 0 1 1 1 0 0 0 1
6 2 0 0 1 0 1 0 0 0 1 0 0 0
4 4 0 0 0 1 1 0 0 1 1 1 0 0 1 0 1 1
5 6 1 0 0 0 0 1 1 1 1 0 0 1 1 0 1 1 1 1 0 0 1 0 0 1 0 1 1 0 0 1 1 0
4 6 0 0 0 1 1 0 1 0 1 1 1 1 1 1 1 1 1 0 1 0 0 1 1 0 0 0
4 5 0 0 1 1 1 1 0 0 0 1 0 0 0 1 1 0 1 0 0 1
3 4 0 0 0 1 1 1 1 0 1 0 1 1 0 0 1
4 5 1 0 0 1 1 1 0 0 0 0 0 0 0 1 0 0 1 1 1 1
5 6 1 0 1 0 0 0 0 1 0 0 0 1 0 1 1 0 0 1 0 0 0 0 1 0 1 0 1 1 1 1 1 1 0 1 1
2 5 1 1 0 0 0 1 0 1 1 0
2 1 0 1 1
3 6 0 1 0 0 0 1 0 1 0 0 1 1 1 0 1 0 1 0 1 1
4 1 0 0 1 0
5 4 1 0 1 1 1 1 0 0 0 0 1 0 1 0 1 0 1 0 1 0
6 4 1 1 1 1 0 0 1 1 0 0 1 1 0 0 1 0 0 0 1 1 1 0 1 0
5 6 1 1 0 0 1 1 1 1 1 1 0 0 0 0 1 0 1 0 0 0 1 0 0 1 1 1 0 0 0 1 0 0
4 6 1 1 0 0 0 1 0 1 0 0 0 1 1 1 0 1 1 0 0 0 0 0 1 1 1 0 0
6 6 0 0 0 0 0 1 1 1 0 0 0 1 1 1 0 0 0 0 0 1 1 1 0 0 0 1 0 0 1 0 0 0 0 1 1 1
5 2 0 0 1 0 0 1 0 0 1 1
5 5 1 0 0 0 0 0 0 1 0 1 1 0 0 1 1 1 1 0 1 0 1 1 1 1 1
6 1 1 1 1 0 0 0 1
3 2 0 1 1 0 1 1 1 0
2 1 0 1 0
2 2 1 0 1 0`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCase(n, m int, grid [][]int) int {
	zeros := make([]int, n+m-1)
	ones := make([]int, n+m-1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 {
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
	return ans
}

type testCase struct {
	n    int
	m    int
	grid [][]int
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("invalid test data")
	}
	idx := 0
	var tests []testCase
	for idx < len(fields) {
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("invalid test data")
		}
		n, _ := strconv.Atoi(fields[idx])
		m, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		total := n * m
		if idx+total > len(fields) {
			return nil, fmt.Errorf("invalid test data")
		}
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				v, _ := strconv.Atoi(fields[idx])
				idx++
				grid[i][j] = v
			}
		}
		tests = append(tests, testCase{n: n, m: m, grid: grid})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := strconv.Itoa(solveCase(tc.n, tc.m, tc.grid))
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
