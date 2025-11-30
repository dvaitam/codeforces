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
	input      string
	expectYes  bool
	expectGrid [][]int
}

const testcaseData = `1 1 1
1 1 0
4 3 1 0 1 1 0 1 0 0 1 0 1 1 0 1 1 1
4 1 0 1 0 0 0 1 1 0 0 1 0 1 1 1 0 0
2 1 1 1 0 1
4 2 0 0 0 0 0 0 1 1 1 1 1 1 0 1 0 0
4 4 0 1 0 0 1 1 1 0 1 1 1 1 0 0 1 1
4 2 0 0 0 1 1 0 1 1 1 1 1 1 1 1 1 1
4 2 1 1 1 1 1 1 0 1 1 1 1 1 0 1 0 1
1 1 0
1 1 1
4 1 0 0 1 0 1 1 0 1 0 0 1 1 0 0 1 0
1 1 0
1 1 1
1 1 0
1 1 1
4 4 0 0 0 1 1 1 1 1 1 0 1 1 1 1 0 0
3 2 1 1 1 1 1 0 0 0 0
1 1 1
3 1 0 0 1 1 0 0 0 1 0
4 4 0 0 0 1 0 0 0 0 1 0 1 0 0 0 1 0
3 1 0 1 0 1 1 1 0 1 0
2 2 1 0 0 0
2 2 0 0 0 0
4 1 0 0 1 0 0 0 1 0 1 0 1 0 0 0 1 0
1 1 0
1 1 1
3 2 0 0 0 1 0 1 1 0 0
3 2 0 0 0 1 1 1 0 1 0
1 1 1
1 1 0
4 3 0 1 0 0 0 1 0 0 1 1 1 1 1 0 1 0
2 2 0 0 0 0
2 1 0 1 1 1
4 2 1 0 1 0 1 1 1 1 0 1 1 1 0 0 1 0
1 1 0
1 1 1
4 2 0 0 1 1 1 1 1 0 1 0 0 0 0 0 0 0
1 1 1
3 1 0 0 0 0 0 1 1 0 0
4 2 1 0 1 1 0 1 0 0 0 1 1 1 0 0 1 0
4 3 0 1 0 0 1 1 1 0 0 0 1 1 0 0 1 1
1 1 0
4 2 0 1 1 1 1 1 1 1 0 0 0 1 1 0 0 1
1 1 1
2 1 1 0 1 1
3 1 0 0 0 1 0 1 0 1 0
4 4 1 0 0 1 0 0 1 1 0 0 1 1 1 1 1 1
3 1 1 0 1 0 1 0 0 0 0
4 4 0 0 1 1 0 0 1 1 1 1 1 0 0 1 1 1
3 2 0 0 0 1 0 0 1 0 0
3 1 0 0 1 1 1 0 0 1 0
2 1 0 0 0 1
4 4 1 1 1 1 0 0 0 0 0 1 0 1 1 1 0 0
1 1 0
4 1 0 1 1 1 1 1 1 1 1 1 0 1 0 1 1 1
2 1 1 0 1 1
4 3 0 1 1 1 0 0 0 1 1 0 0 1 1 0 0 1
1 1 0
1 1 0
3 3 1 1 1 0 1 1 1 0 1
3 1 0 0 0 0 0 1 0 1 1
3 2 0 0 1 0 0 0 0 0 1
4 2 1 1 1 0 1 1 0 0 0 0 1 0 1 0 0 0
1 1 0
4 3 0 1 0 0 0 1 0 1 0 0 0 1 0 1 0 0
4 4 0 1 0 0 0 1 1 1 1 0 1 0 0 1 0 1
3 1 1 1 1 1 1 1 1 1 0
3 3 0 1 0 0 0 0 0 1 1
1 1 1
1 1 1
4 1 1 1 0 1 0 1 0 1 1 0 1 0 0 0 1 0
2 1 1 0 1 1
2 2 1 0 0 0
2 1 1 0 1 0
3 2 1 0 0 0 0 1 0 0 0
1 1 1
4 4 1 0 1 1 1 1 1 0 0 0 1 1 0 1 1 0
1 1 0
2 1 1 0 0 1
2 1 0 1 0 1
4 1 0 1 0 0 0 1 1 0 0 1 1 1 0 1 0 0
4 4 0 0 1 1 0 1 0 1 1 1 1 1 1 1 1 1
2 1 1 1 0 0
3 3 0 0 0 0 0 1 0 1 0
2 1 0 0 1 0
2 1 0 0 0 1
4 1 0 1 0 1 0 1 1 0 1 1 1 0 1 1 0 0
2 2 0 0 0 1
3 1 0 1 1 1 0 0 1 0 1
4 1 1 1 0 0 1 1 0 0 0 0 0 0 1 0 1 0
2 2 1 1 0 0
2 2 1 0 0 0
2 2 0 1 1 1
3 1 0 0 1 0 0 0 1 0 1
3 3 1 0 1 0 0 1 0 0 0
2 1 0 1 0 1
1 1 1
1 1 1
1 1 1`

func solveG(n, k int, grid [][]int) (bool, [][]int) {
	f := make([][]int, n+2)
	vst := make([][]bool, n+2)
	for i := 0; i < n+2; i++ {
		f[i] = make([]int, n+2)
		vst[i] = make([]bool, n+2)
	}
	mx := make([][]int, k)
	for i := 0; i < k; i++ {
		mx[i] = make([]int, n+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if grid[i-1][j-1] == 0 {
				f[i][j] = 1
			} else {
				f[i][j] = 0
			}
		}
	}
	isNo := false
	for i := n; i >= 1 && !isNo; i-- {
		for j := n; j >= 1; j-- {
			if f[i+1][j+1] > 0 {
				f[i][j] += f[i+1][j+1]
			}
			if f[i+1][j] > f[i][j] {
				f[i][j] = f[i+1][j]
			}
			if f[i][j+1] > f[i][j] {
				f[i][j] = f[i][j+1]
			}
			if f[i][j] == k {
				isNo = true
				break
			}
			if mx[f[i][j]][j] == 0 {
				mx[f[i][j]][j] = i
			}
		}
	}
	if isNo {
		return false, nil
	}
	for level := k - 1; level >= 1; level-- {
		for j := n - 1; j >= 1; j-- {
			if mx[level][j+1] > mx[level][j] {
				mx[level][j] = mx[level][j+1]
			}
		}
		x, y := n, 1
		for y <= n && vst[x][y] {
			y++
		}
		for y <= n {
			vst[x][y] = true
			if (y == n || x != mx[level][y+1]) && x > 1 && !vst[x-1][y] {
				x--
			} else {
				y++
			}
		}
	}
	res := make([][]int, n)
	for i := range res {
		res[i] = make([]int, n)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if vst[i][j] {
				res[i-1][j-1] = 1
			}
		}
	}
	return true, res
}

func parseTestcase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return testCase{}, fmt.Errorf("invalid testcase")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return testCase{}, fmt.Errorf("bad k: %w", err)
	}
	if len(fields) != 2+n*n {
		return testCase{}, fmt.Errorf("expected %d cells, got %d", n*n, len(fields)-2)
	}
	grid := make([][]int, n)
	idx := 2
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[idx])
			if err != nil {
				return testCase{}, fmt.Errorf("bad cell %d,%d: %w", i, j, err)
			}
			grid[i][j] = val
			idx++
		}
	}

	ok, expectedGrid := solveG(n, k, grid)
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			input.WriteByte(byte('0' + grid[i][j]))
		}
		input.WriteByte('\n')
	}

	return testCase{input: input.String(), expectYes: ok, expectGrid: expectedGrid}, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseTestcase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierG /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
			fmt.Printf("test %d: empty output\n", idx+1)
			os.Exit(1)
		}
		first := strings.ToUpper(strings.TrimSpace(lines[0]))
		if !tc.expectYes {
			if first != "NO" {
				fmt.Printf("test %d failed: expected NO got %s\n", idx+1, lines[0])
				os.Exit(1)
			}
			continue
		}
		if first != "YES" {
			fmt.Printf("test %d failed: expected YES got %s\n", idx+1, lines[0])
			os.Exit(1)
		}
		if len(lines)-1 != len(tc.expectGrid) {
			fmt.Printf("test %d failed: expected %d grid lines got %d\n", idx+1, len(tc.expectGrid), len(lines)-1)
			os.Exit(1)
		}
		for i := 0; i < len(tc.expectGrid); i++ {
			if len(lines[i+1]) != len(tc.expectGrid[i]) {
				fmt.Printf("test %d line %d length mismatch\n", idx+1, i+1)
				os.Exit(1)
			}
			for j := 0; j < len(tc.expectGrid[i]); j++ {
				if int(lines[i+1][j]-'0') != tc.expectGrid[i][j] {
					fmt.Printf("test %d grid mismatch\n", idx+1)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
