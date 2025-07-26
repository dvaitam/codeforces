package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	n, m     int
	grid     []string
	solvable bool
}

func parseTestCases(path string) ([]testCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var cases []testCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) < 3 {
			return nil, fmt.Errorf("bad testcase line: %s", line)
		}
		n, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("bad n: %v", err)
		}
		m, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("bad m: %v", err)
		}
		if len(tokens) != 2+n {
			return nil, fmt.Errorf("expected %d rows, got %d", n, len(tokens)-2)
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := tokens[2+i]
			if len(row) != m {
				return nil, fmt.Errorf("row %d length mismatch", i)
			}
			grid[i] = row
		}
		// compute solvable
		dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		solvable := true
		for i := 0; i < n && solvable; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'W' {
					for _, d := range dirs {
						ni, nj := i+d[0], j+d[1]
						if ni >= 0 && ni < n && nj >= 0 && nj < m {
							if grid[ni][nj] == 'S' {
								solvable = false
								break
							}
						}
					}
				}
			}
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, m)
		for _, row := range grid {
			buf.WriteString(row)
			buf.WriteByte('\n')
		}
		cases = append(cases, testCase{input: buf.String(), n: n, m: m, grid: grid, solvable: solvable})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return string(out), nil
}

func verifyOutput(out string, tc testCase) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(lines[0])
	if first == "No" || first == "NO" || first == "no" {
		if tc.solvable {
			return fmt.Errorf("expected Yes, got No")
		}
		return nil
	}
	if !(first == "Yes" || first == "YES" || first == "yes") {
		return fmt.Errorf("first line should be Yes or No")
	}
	if !tc.solvable {
		return fmt.Errorf("expected No, got Yes")
	}
	if len(lines)-1 != tc.n {
		return fmt.Errorf("expected %d grid lines, got %d", tc.n, len(lines)-1)
	}
	outGrid := make([][]byte, tc.n)
	for i := 0; i < tc.n; i++ {
		row := strings.TrimSpace(lines[i+1])
		if len(row) != tc.m {
			return fmt.Errorf("row %d length mismatch", i)
		}
		outGrid[i] = []byte(row)
		for j := 0; j < tc.m; j++ {
			ch := outGrid[i][j]
			inCh := tc.grid[i][j]
			if inCh == 'S' || inCh == 'W' {
				if ch != inCh {
					return fmt.Errorf("cell %d,%d changed from %c to %c", i, j, inCh, ch)
				}
			} else {
				if ch != '.' && ch != 'D' {
					return fmt.Errorf("invalid char %c at %d,%d", ch, i, j)
				}
			}
		}
	}
	dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if outGrid[i][j] == 'W' {
				for _, d := range dirs {
					ni, nj := i+d[0], j+d[1]
					if ni >= 0 && ni < tc.n && nj >= 0 && nj < tc.m {
						if outGrid[ni][nj] == 'S' {
							return fmt.Errorf("wolf at %d,%d adjacent to sheep", i, j)
						}
					}
				}
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cases, err := parseTestCases("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
