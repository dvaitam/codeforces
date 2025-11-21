package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	n     int
	m     int
	grid  [][]byte
}

type operation struct {
	x, y, z int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	if out == "-1" {
		if possible(tc.grid) {
			return fmt.Errorf("reported -1 but redistribution is possible")
		}
		return nil
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read number of operations: %v", err)
	}
	ops := make([]operation, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &ops[i].x, &ops[i].y, &ops[i].z); err != nil {
			return fmt.Errorf("failed to read operation %d: %v", i+1, err)
		}
	}
	return validateOperations(tc, ops)
}

func validateOperations(tc testCase, ops []operation) error {
	n, m := tc.n, tc.m
	grid := copyGrid(tc.grid)
	target := totalOnes(grid) / n
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				counts[i]++
			}
		}
	}
	for idx, op := range ops {
		x, y, z := op.x-1, op.y-1, op.z-1
		if x < 0 || x >= n || y < 0 || y >= n || z < 0 || z >= m {
			return fmt.Errorf("operation %d out of bounds", idx+1)
		}
		if grid[x][z] == grid[y][z] {
			return fmt.Errorf("operation %d swaps identical values", idx+1)
		}
		if grid[x][z] == '1' && counts[x] <= target {
			return fmt.Errorf("operation %d moves 1 from non-surplus row %d", idx+1, x+1)
		}
		if grid[y][z] == '1' && counts[y] <= target {
			return fmt.Errorf("operation %d moves 1 from non-surplus row %d", idx+1, y+1)
		}
		if grid[x][z] == '1' && counts[y] >= target && grid[y][z] == '0' {
			return fmt.Errorf("operation %d adds 1 to non-deficit row %d", idx+1, y+1)
		}
		if grid[y][z] == '1' && counts[x] >= target && grid[x][z] == '0' {
			return fmt.Errorf("operation %d adds 1 to non-deficit row %d", idx+1, x+1)
		}
		grid[x][z], grid[y][z] = grid[y][z], grid[x][z]
		counts[x] = rowCount(grid[x])
		counts[y] = rowCount(grid[y])
	}
	for i := 0; i < n; i++ {
		if rowCount(grid[i]) != target {
			return fmt.Errorf("row %d has %d ones instead of %d", i+1, rowCount(grid[i]), target)
		}
	}
	return nil
}

func possible(grid [][]byte) bool {
	n := len(grid)
	total := totalOnes(grid)
	return total%n == 0
}

func totalOnes(grid [][]byte) int {
	total := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == '1' {
				total++
			}
		}
	}
	return total
}

func rowCount(row []byte) int {
	cnt := 0
	for _, v := range row {
		if v == '1' {
			cnt++
		}
	}
	return cnt
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTestCase([]string{"10", "01"}),
		makeTestCase([]string{"111", "000", "000"}),
	}
	for i := 0; i < 50; i++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if rand.Intn(2) == 0 {
					row[c] = '0'
				} else {
					row[c] = '1'
				}
			}
			grid[r] = string(row)
		}
		tests = append(tests, makeTestCase(grid))
	}
	return tests
}

func makeTestCase(rows []string) testCase {
	n, m := len(rows), len(rows[0])
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	for _, row := range rows {
		fmt.Fprintln(&sb, row)
	}
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = []byte(rows[i])
	}
	return testCase{
		input: sb.String(),
		n:     n,
		m:     m,
		grid:  grid,
	}
}

func copyGrid(grid [][]byte) [][]byte {
	n := len(grid)
	cp := make([][]byte, n)
	for i := 0; i < n; i++ {
		cp[i] = append([]byte(nil), grid[i]...)
	}
	return cp
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
