package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
}

var testcases = []testCase{
	{input: "3 2 1 1 0 1 0 1 1 1 0 1 0 1"},
	{input: "3 3 0 01 11 10 00 11 01 0 0 0 01 11 10 00 11 01 0 0"},
	{input: "3 3 1 10 00 0 0 11 1 0 00 1 10 00 0 0 11 1 0 00"},
	{input: "2 3 0 01 10 00 0 1 0 01 10 00 0 1"},
	{input: "3 3 01 0 0 00 0 01 1 0 01 01 0 0 00 0 01 1 0 01"},
	{input: "3 3 10 1 11 10 0 10 0 0 00 10 1 11 10 0 10 0 0 00"},
	{input: "3 3 11 11 0 11 0 1 1 00 0 11 11 0 11 0 1 1 00 0"},
	{input: "3 2 01 11 00 10 1 1 01 11 00 10 1 1"},
	{input: "3 3 0 01 11 1 01 10 1 1 01 0 01 11 1 01 10 1 1 01"},
	{input: "2 3 0 00 10 1 0 01 0 00 10 1 0 01"},
	{input: "3 3 1 10 0 0 01 11 10 0 01 1 10 0 0 01 11 10 0 01"},
	{input: "3 3 11 0 0 1 0 0 1 1 0 11 0 0 1 0 0 1 1 0"},
	{input: "3 3 0 0 0 11 1 0 1 0 1 0 0 0 11 1 0 1 0 1"},
	{input: "3 2 01 0 0 1 0 01 01 0 0 1 0 01"},
	{input: "2 2 11 1 0 11 11 1 0 11"},
	{input: "3 3 1 00 00 1 1 00 1 11 1 1 00 00 1 1 00 1 11 1"},
	{input: "3 3 0 1 10 01 0 01 0 0 0 0 1 10 01 0 01 0 0 0"},
	{input: "3 3 0 0 00 0 1 11 1 00 1 0 0 00 0 1 11 1 00 1"},
	{input: "2 2 0 1 01 11 0 1 01 11"},
	{input: "3 2 1 00 1 0 10 00 1 00 1 0 10 00"},
	{input: "3 3 1 0 11 00 01 0 10 1 00 1 0 11 00 01 0 10 1 00"},
	{input: "2 2 10 1 01 0 10 1 01 0"},
	{input: "3 2 0 1 11 0 0 00 0 1 11 0 0 00"},
	{input: "2 2 0 0 1 1 0 0 1 1"},
	{input: "3 3 00 1 1 11 01 0 0 0 10 00 1 1 11 01 0 0 0 10"},
	{input: "3 3 1 01 1 01 01 1 1 01 10 1 01 1 01 01 1 1 01 10"},
	{input: "2 2 01 01 11 1 01 01 11 1"},
	{input: "3 2 0 10 1 10 10 0 0 10 1 10 10 0"},
	{input: "3 3 01 11 00 00 1 01 1 1 1 01 11 00 00 1 01 1 1 1"},
	{input: "3 2 00 0 11 0 11 1 00 0 11 0 11 1"},
	{input: "3 2 0 1 1 1 0 00 0 1 1 1 0 00"},
	{input: "2 2 00 10 01 0 00 10 01 0"},
	{input: "2 3 1 1 01 01 1 1 1 1 01 01 1 1"},
	{input: "3 3 11 10 0 10 11 0 1 0 10 11 10 0 10 11 0 1 0 10"},
	{input: "2 3 11 00 0 00 10 1 11 00 0 00 10 1"},
	{input: "2 3 00 11 00 1 11 1 00 11 00 1 11 1"},
	{input: "3 3 1 10 11 11 10 01 11 0 01 1 10 11 11 10 01 11 0 01"},
	{input: "2 3 0 00 0 0 1 01 0 00 0 0 1 01"},
	{input: "3 3 00 00 1 11 01 10 1 1 01 00 00 1 11 01 10 1 1 01"},
	{input: "2 3 0 1 1 00 0 0 0 1 1 00 0 0"},
	{input: "3 2 0 0 0 0 1 00 0 0 0 0 1 00"},
	{input: "3 3 0 1 00 0 1 11 00 1 11 0 1 00 0 1 11 00 1 11"},
	{input: "3 2 0 1 0 0 10 10 0 1 0 0 10 10"},
	{input: "2 2 10 1 1 00 10 1 1 00"},
	{input: "2 3 01 10 01 1 01 1 01 10 01 1 01 1"},
	{input: "2 3 01 0 10 1 0 10 01 0 10 1 0 10"},
	{input: "3 2 0 01 0 0 11 10 0 01 0 0 11 10"},
	{input: "2 3 01 1 0 11 0 11 01 1 0 11 0 11"},
	{input: "2 3 0 1 0 1 1 11 0 1 0 1 1 11"},
	{input: "3 3 01 11 10 01 10 1 1 0 11 01 11 10 01 10 1 1 0 11"},
	{input: "3 3 00 01 00 0 01 10 11 00 11 00 01 00 0 01 10 11 00 11"},
	{input: "2 3 0 0 00 1 00 1 0 0 00 1 00 1"},
	{input: "2 3 1 01 01 0 1 10 1 01 01 0 1 10"},
	{input: "2 3 10 1 11 10 01 1 10 1 11 10 01 1"},
	{input: "2 2 11 11 0 01 11 11 0 01"},
	{input: "2 3 0 10 1 1 00 11 0 10 1 1 00 11"},
	{input: "3 3 0 00 1 0 00 01 10 0 0 0 00 1 0 00 01 10 0 0"},
	{input: "2 3 1 0 1 01 10 10 1 0 1 01 10 10"},
	{input: "2 2 1 01 11 0 1 01 11 0"},
	{input: "3 2 00 0 10 00 00 1 00 0 10 00 00 1"},
	{input: "3 3 11 0 00 10 0 11 10 10 1 11 0 00 10 0 11 10 10 1"},
	{input: "2 2 0 0 1 00 0 0 1 00"},
	{input: "2 2 1 11 1 0 1 11 1 0"},
	{input: "2 3 1 0 0 11 11 11 1 0 0 11 11 11"},
	{input: "3 3 11 10 1 11 0 00 0 00 01 11 10 1 11 0 00 0 00 01"},
	{input: "2 3 0 0 0 0 1 00 0 0 0 0 1 00"},
	{input: "3 3 1 0 10 1 0 0 11 1 00 1 0 10 1 0 0 11 1 00"},
	{input: "3 2 0 10 10 01 10 01 0 10 10 01 10 01"},
	{input: "3 2 1 00 1 01 11 1 1 00 1 01 11 1"},
	{input: "2 3 0 11 1 01 1 1 0 11 1 01 1 1"},
	{input: "2 2 1 01 01 0 1 01 01 0"},
	{input: "2 2 0 00 11 11 0 00 11 11"},
	{input: "2 3 10 01 01 00 00 10 10 01 01 00 00 10"},
	{input: "3 3 1 0 1 00 0 1 00 0 10 1 0 1 00 0 1 00 0 10"},
	{input: "2 3 0 0 1 1 0 01 0 0 1 1 0 01"},
	{input: "3 2 01 01 10 01 0 01 01 01 10 01 0 01"},
	{input: "2 3 10 1 0 0 1 00 10 1 0 0 1 00"},
	{input: "2 2 1 10 01 01 1 10 01 01"},
	{input: "2 2 10 1 00 0 10 1 00 0"},
	{input: "2 2 11 00 10 01 11 00 10 01"},
	{input: "2 3 1 01 00 1 00 0 1 01 00 1 00 0"},
	{input: "2 3 1 0 1 01 0 1 1 0 1 01 0 1"},
	{input: "3 3 0 01 1 11 1 1 01 0 0 0 01 1 11 1 1 01 0 0"},
	{input: "3 2 01 00 01 1 0 10 01 00 01 1 0 10"},
	{input: "2 2 10 10 0 01 10 10 0 01"},
	{input: "3 3 00 10 1 11 10 01 1 0 1 00 10 1 11 10 01 1 0 1"},
	{input: "2 3 1 00 11 0 0 0 1 00 11 0 0 0"},
	{input: "3 3 01 0 00 00 00 1 11 0 1 01 0 00 00 00 1 11 0 1"},
	{input: "3 3 1 01 0 1 11 01 01 11 0 1 01 0 1 11 01 01 11 0"},
	{input: "3 3 00 01 0 0 01 0 1 10 1 00 01 0 0 01 0 1 10 1"},
	{input: "3 3 0 00 0 0 00 1 11 10 0 0 00 0 0 00 1 11 10 0"},
	{input: "2 2 0 0 1 11 0 0 1 11"},
	{input: "2 3 0 11 1 1 10 11 0 11 1 1 10 11"},
	{input: "3 3 1 01 00 00 0 1 00 01 10 1 01 00 00 0 1 00 01 10"},
	{input: "3 3 11 1 1 0 11 10 1 11 0 11 1 1 0 11 10 1 11 0"},
	{input: "2 2 11 01 01 01 11 01 01 01"},
	{input: "3 3 01 1 1 11 1 0 00 1 10 01 1 1 11 1 0 00 1 10"},
	{input: "3 2 1 00 01 1 0 0 1 00 01 1 0 0"},
	{input: "3 3 01 11 00 1 1 00 0 1 0 01 11 00 1 1 00 0 1 0"},
	{input: "3 3 01 00 11 01 0 0 00 0 0 01 00 11 01 0 0 00 0 0"},
}

func parseGrid(scanner *bufio.Scanner, n, m int) [][]string {
	grid := make([][]string, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]string, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			grid[i][j] = scanner.Text()
		}
	}
	return grid
}

func parseInput(input string) (int, int, [][]string, [][]string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	n := 0
	fmt.Sscan(scanner.Text(), &n)
	scanner.Scan()
	m := 0
	fmt.Sscan(scanner.Text(), &m)
	init := parseGrid(scanner, n, m)
	fin := parseGrid(scanner, n, m)
	return n, m, init, fin
}

func simulateOps(n, m int, grid [][]string, ops [][4]int) error {
	for idx, op := range ops {
		x1, y1, x2, y2 := op[0]-1, op[1]-1, op[2]-1, op[3]-1
		if x1 < 0 || x1 >= n || y1 < 0 || y1 >= m || x2 < 0 || x2 >= n || y2 < 0 || y2 >= m {
			return fmt.Errorf("op %d: coords out of range", idx+1)
		}
		if x1 == x2 && y1 == y2 {
			return fmt.Errorf("op %d: same cell", idx+1)
		}
		if x1 != x2 && y1 != y2 {
			return fmt.Errorf("op %d: not same row or column", idx+1)
		}
		if len(grid[x1][y1]) == 0 {
			return fmt.Errorf("op %d: source cell empty", idx+1)
		}
		s := grid[x1][y1]
		ch := s[len(s)-1]
		grid[x1][y1] = s[:len(s)-1]
		grid[x2][y2] = string(ch) + grid[x2][y2]
	}
	return nil
}

func runCase(exe, input string) error {
	n, m, initGrid, finGrid := parseInput(input)

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	// Parse output
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out.String())))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	var q int
	fmt.Sscan(scanner.Text(), &q)

	// Check q limit
	totalS := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			totalS += len(initGrid[i][j])
		}
	}
	if q < 0 || q > 4*totalS {
		return fmt.Errorf("q=%d exceeds 4*s=%d", q, 4*totalS)
	}

	ops := make([][4]int, q)
	for i := 0; i < q; i++ {
		for j := 0; j < 4; j++ {
			if !scanner.Scan() {
				return fmt.Errorf("missing op data at op %d", i+1)
			}
			fmt.Sscan(scanner.Text(), &ops[i][j])
		}
	}

	if err := simulateOps(n, m, initGrid, ops); err != nil {
		return err
	}

	// Check final state matches
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if initGrid[i][j] != finGrid[i][j] {
				return fmt.Errorf("cell (%d,%d) is %q, expected %q", i+1, j+1, initGrid[i][j], finGrid[i][j])
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := tc.input + "\n"
		if err := runCase(bin, input); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
