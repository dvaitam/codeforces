package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in         string
	a, b, c, d int
}

func generate() []testCase {
	const T = 100
	rand.Seed(3)
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		a := rand.Intn(100) + 1
		b := rand.Intn(100) + 1
		c := rand.Intn(100) + 1
		d := rand.Intn(100) + 1
		cases[i] = testCase{
			in: fmt.Sprintf("%d %d %d %d\n", a, b, c, d),
			a:  a,
			b:  b,
			c:  c,
			d:  d,
		}
	}
	return cases
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
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func check(tc testCase, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	var n, m int
	if _, err := fmt.Sscanf(strings.TrimSpace(lines[0]), "%d %d", &n, &m); err != nil {
		return fmt.Errorf("failed to parse dimensions: %v", err)
	}
	if n < 1 || n > 50 || m < 1 || m > 50 {
		return fmt.Errorf("invalid dimensions %d %d", n, m)
	}
	if len(lines) != n+1 {
		return fmt.Errorf("expected %d lines, got %d", n+1, len(lines))
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		line := strings.TrimSpace(lines[i+1])
		if len(line) != m {
			return fmt.Errorf("line %d length %d, expected %d", i+1, len(line), m)
		}
		row := []byte(line)
		for j := 0; j < m; j++ {
			ch := row[j]
			if ch < 'A' || ch > 'D' {
				return fmt.Errorf("invalid char %c at (%d,%d)", ch, i, j)
			}
		}
		grid[i] = row
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}
	count := map[byte]int{'A': 0, 'B': 0, 'C': 0, 'D': 0}
	queue := make([][2]int, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			ch := grid[i][j]
			count[ch]++
			queue = queue[:0]
			queue = append(queue, [2]int{i, j})
			visited[i][j] = true
			for len(queue) > 0 {
				p := queue[0]
				queue = queue[1:]
				for _, d := range dirs {
					ni, nj := p[0]+d[0], p[1]+d[1]
					if ni >= 0 && ni < n && nj >= 0 && nj < m && !visited[ni][nj] && grid[ni][nj] == ch {
						visited[ni][nj] = true
						queue = append(queue, [2]int{ni, nj})
					}
				}
			}
		}
	}
	if count['A'] != tc.a || count['B'] != tc.b || count['C'] != tc.c || count['D'] != tc.d {
		return fmt.Errorf("component counts mismatch: got A:%d B:%d C:%d D:%d", count['A'], count['B'], count['C'], count['D'])
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		out, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := check(tc, out); err != nil {
			fmt.Printf("case %d failed\ninput:\n%s\nerror: %v\noutput:\n%s\n", idx+1, tc.in, err, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
