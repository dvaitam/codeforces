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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
		expectedLines := solveB(tc.n, tc.grid)
		expected := strings.Join(expectedLines, "\n")
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expected {
			fmt.Printf("case %d failed: expected:\n%s\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
