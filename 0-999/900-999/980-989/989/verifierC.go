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
	if len(lines) != 51 {
		return fmt.Errorf("expected 51 lines, got %d", len(lines))
	}
	if strings.TrimSpace(lines[0]) != "50 50" {
		return fmt.Errorf("first line should be '50 50', got '%s'", lines[0])
	}
	grid := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		line := strings.TrimSpace(lines[i+1])
		if len(line) != 50 {
			return fmt.Errorf("line %d length %d, expected 50", i+1, len(line))
		}
		row := []byte(line)
		for j := 0; j < 50; j++ {
			ch := row[j]
			if ch < 'A' || ch > 'D' {
				return fmt.Errorf("invalid char %c at (%d,%d)", ch, i, j)
			}
		}
		grid[i] = row
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	visited := make([][]bool, 50)
	for i := 0; i < 50; i++ {
		visited[i] = make([]bool, 50)
	}
	count := map[byte]int{'A': 0, 'B': 0, 'C': 0, 'D': 0}
	queue := make([][2]int, 0, 2500)
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
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
					if ni >= 0 && ni < 50 && nj >= 0 && nj < 50 && !visited[ni][nj] && grid[ni][nj] == ch {
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
