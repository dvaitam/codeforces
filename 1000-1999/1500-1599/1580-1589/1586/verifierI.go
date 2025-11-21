package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type testCase struct {
	input     string
	n         int
	given     []string
	fullCount int
	solutions [][]string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		buildCase(2, []string{"..", ".."}),
		buildCase(2, []string{"SS", "SS"}),
		buildCase(4, []string{"S..G", "....", "....", "G..S"}),
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomCase(4))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomCase(6))
	}
	tests = append(tests, randomCase(8))
	return tests
}

func randomCase(n int) testCase {
	given := make([]string, n)
	for i := range given {
		row := make([]byte, n)
		for j := range row {
			v := rand.Intn(4)
			switch v {
			case 0:
				row[j] = '.'
			case 1:
				row[j] = 'S'
			case 2:
				row[j] = 'G'
			default:
				row[j] = '.'
			}
		}
		given[i] = string(row)
	}
	return buildCase(n, given)
}

func buildCase(n int, given []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(given[i])
		sb.WriteByte('\n')
	}
	solutions := enumerateSolutions(n, given)
	return testCase{
		input:     sb.String(),
		n:         n,
		given:     given,
		fullCount: len(solutions),
		solutions: solutions,
	}
}

func enumerateSolutions(n int, given []string) [][]string {
	if n%2 == 1 {
		return nil
	}
	var sols [][]string
	for mask := 0; mask < 4; mask++ {
		var grid []string
		switch mask {
		case 0:
			grid = generateStripes(n, false)
		case 1:
			grid = generateStripes(n, true)
		case 2:
			grid = generateRings(n, false)
		case 3:
			grid = generateRings(n, true)
		}
		if matchesGrid(grid, given) {
			sols = append(sols, grid)
		}
	}
	return sols
}

func generateStripes(n int, swap bool) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			val := (i/2 + j/2) % 2
			if swap {
				val ^= 1
			}
			if val == 0 {
				row[j] = 'S'
			} else {
				row[j] = 'G'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func generateRings(n int, swap bool) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			layer := min(i, j, n-1-i, n-1-j)
			val := layer % 2
			if swap {
				val ^= 1
			}
			if val == 0 {
				row[j] = 'S'
			} else {
				row[j] = 'G'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func matchesGrid(grid []string, given []string) bool {
	n := len(grid)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if given[i][j] != '.' && given[i][j] != grid[i][j] {
				return false
			}
		}
	}
	return true
}

func checkOutput(out string, tc testCase) error {
	if tc.n == 0 {
		return fmt.Errorf("invalid test case size")
	}
	solutions := tc.solutions
	if tc.n%2 == 1 {
		if out != "NONE" {
			return fmt.Errorf("expected NONE for odd n, got %q", out)
		}
		return nil
	}
	switch tc.fullCount {
	case 0:
		if out != "NONE" {
			return fmt.Errorf("expected NONE, got %q", out)
		}
	case 1:
		lines := strings.Split(out, "\n")
		if len(lines) < tc.n+1 {
			return fmt.Errorf("output too short; expected header + %d rows", tc.n)
		}
		if lines[0] != "UNIQUE" {
			return fmt.Errorf("expected UNIQUE, got %q", lines[0])
		}
		grid := lines[1 : tc.n+1]
		if err := verifyGrid(grid, tc.given); err != nil {
			return err
		}
		ref := solutions[0]
		for i := 0; i < tc.n; i++ {
			if grid[i] != ref[i] {
				return fmt.Errorf("unique solution mismatch at row %d", i)
			}
		}
	default:
		if out == "NONE" || out == "UNIQUE" {
			return fmt.Errorf("expected MULTIPLE, got %q", out)
		}
		if out == "MULTIPLE" {
			return nil
		}
		lines := strings.Split(out, "\n")
		if len(lines) < tc.n+1 {
			return fmt.Errorf("output too short for grid")
		}
		header := lines[0]
		if header != "MULTIPLE" {
			return fmt.Errorf("expected MULTIPLE header, got %q", header)
		}
		grid := lines[1 : tc.n+1]
		if err := verifyGrid(grid, tc.given); err != nil {
			return err
		}
		valid := false
		for _, sol := range solutions {
			match := true
			for i := 0; i < tc.n; i++ {
				if grid[i] != sol[i] {
					match = false
					break
				}
			}
			if match {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New("grid not in valid solution set")
		}
	}
	return nil
}

func verifyGrid(grid []string, given []string) error {
	n := len(given)
	if len(grid) != n {
		return fmt.Errorf("grid row count mismatch: expected %d, got %d", n, len(grid))
	}
	for i := 0; i < n; i++ {
		if len(grid[i]) != n {
			return fmt.Errorf("row %d length mismatch", i+1)
		}
		for j := 0; j < n; j++ {
			ch := grid[i][j]
			if ch != 'S' && ch != 'G' {
				return fmt.Errorf("invalid character %q at (%d,%d)", ch, i+1, j+1)
			}
			if given[i][j] != '.' && given[i][j] != ch {
				return fmt.Errorf("grid conflicts with input at (%d,%d)", i+1, j+1)
			}
			if !hasTwoSameNeighbors(grid, i, j) {
				return fmt.Errorf("cell (%d,%d) does not have two same-colored neighbors", i+1, j+1)
			}
		}
	}
	return nil
}

func hasTwoSameNeighbors(grid []string, i, j int) bool {
	n := len(grid)
	ch := grid[i][j]
	count := 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, d := range dirs {
		ni, nj := i+d[0], j+d[1]
		if ni >= 0 && ni < n && nj >= 0 && nj < n && grid[ni][nj] == ch {
			count++
		}
	}
	return count == 2
}

func min(a ...int) int {
	res := a[0]
	for _, v := range a[1:] {
		if v < res {
			res = v
		}
	}
	return res
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
