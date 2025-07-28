package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Fenwick struct{ tree []int }

func NewFenwick(n int) *Fenwick { return &Fenwick{make([]int, n+2)} }
func (f *Fenwick) Add(i, delta int) {
	for ; i < len(f.tree); i += i & -i {
		f.tree[i] += delta
	}
}
func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m int, grid []string, queries [][2]int) string {
	total := n * m
	fw := NewFenwick(total)
	state := make([]byte, total+1)
	stars := 0
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if grid[i][j] == '*' {
				pos := j*n + i + 1
				state[pos] = 1
				fw.Add(pos, 1)
				stars++
			}
		}
	}
	var out strings.Builder
	for _, q := range queries {
		x, y := q[0]-1, q[1]-1
		pos := y*n + x + 1
		if state[pos] == 1 {
			state[pos] = 0
			fw.Add(pos, -1)
			stars--
		} else {
			state[pos] = 1
			fw.Add(pos, 1)
			stars++
		}
		good := fw.Sum(stars)
		out.WriteString(fmt.Sprintf("%d\n", stars-good))
	}
	return strings.TrimSpace(out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(47))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		q := rng.Intn(8) + 1
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(2) == 0 {
					row[j] = '.'
				} else {
					row[j] = '*'
				}
			}
			grid[i] = string(row)
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			queries[i][0] = rng.Intn(n) + 1
			queries[i][1] = rng.Intn(m) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i] + "\n")
		}
		for i := 0; i < q; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", queries[i][0], queries[i][1]))
		}
		input := sb.String()
		exp := expected(n, m, grid, queries)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
