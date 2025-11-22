package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	n    int
	grid []string
}

func parseTestcasesB(path string) ([]testCaseB, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseB
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("expected %d grid rows", n)
		}
		grid := make([]string, n)
		copy(grid, fields[1:])
		cases = append(cases, testCaseB{n: n, grid: grid})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveB(n int, grid []string) []string {
	g := make([][]byte, n)
	stars := [][2]int{}
	for i := 0; i < n; i++ {
		g[i] = []byte(grid[i])
		for j := 0; j < n; j++ {
			if g[i][j] == '*' {
				stars = append(stars, [2]int{i, j})
			}
		}
	}
	r1, c1 := stars[0][0], stars[0][1]
	r2, c2 := stars[1][0], stars[1][1]
	if r1 == r2 {
		r3 := r1 + 1
		if r3 >= n {
			r3 = r1 - 1
		}
		g[r3][c1] = '*'
		g[r3][c2] = '*'
	} else if c1 == c2 {
		c3 := c1 + 1
		if c3 >= n {
			c3 = c1 - 1
		}
		g[r1][c3] = '*'
		g[r2][c3] = '*'
	} else {
		g[r1][c2] = '*'
		g[r2][c1] = '*'
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(g[i])
	}
	return res
}

func run(bin, input string) (string, error) {
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

func validateCase(tc testCaseB, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != tc.n {
		return fmt.Errorf("expected %d lines, got %d", tc.n, len(lines))
	}
	origStars := make(map[[2]int]struct{}, 2)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if tc.grid[i][j] == '*' {
				origStars[[2]int{i, j}] = struct{}{}
			}
		}
	}

	stars := make(map[[2]int]struct{})
	rowSet := make(map[int]struct{})
	colSet := make(map[int]struct{})
	for i := 0; i < tc.n; i++ {
		if len(lines[i]) != tc.n {
			return fmt.Errorf("line %d has length %d, expected %d", i+1, len(lines[i]), tc.n)
		}
		for j := 0; j < tc.n; j++ {
			if lines[i][j] == '*' {
				stars[[2]int{i, j}] = struct{}{}
				rowSet[i] = struct{}{}
				colSet[j] = struct{}{}
			}
		}
	}

	if len(stars) != 4 {
		return fmt.Errorf("expected 4 stars, got %d", len(stars))
	}
	for pos := range origStars {
		if _, ok := stars[pos]; !ok {
			return fmt.Errorf("missing original star at %v", pos)
		}
	}
	if len(rowSet) != 2 || len(colSet) != 2 {
		return fmt.Errorf("stars do not form an axis-aligned rectangle")
	}
	var rows, cols [2]int
	idx := 0
	for r := range rowSet {
		rows[idx] = r
		idx++
	}
	idx = 0
	for c := range colSet {
		cols[idx] = c
		idx++
	}
	corners := [][2]int{{rows[0], cols[0]}, {rows[0], cols[1]}, {rows[1], cols[0]}, {rows[1], cols[1]}}
	for _, p := range corners {
		if _, ok := stars[p]; !ok {
			return fmt.Errorf("missing corner star at %v", p)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		if len(os.Args) < 3 {
			fmt.Println("usage: go run verifierB.go /path/to/binary")
			os.Exit(1)
		}
		bin = os.Args[2]
	}
	cases, err := parseTestcasesB("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(tc.grid[i])
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validateCase(tc, got); err != nil {
			fmt.Printf("case %d failed: %v\noutput:\n%s\n", idx+1, err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
