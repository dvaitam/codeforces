package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	in  string
	out string
}

func compute(grid [][]byte) int {
	n := len(grid)
	m := len(grid[0])
	size := n * m
	dist := make([]int, size)
	queue := make([]int, 0, size)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			u := grid[i][j]
			if u == '<' || u == '>' || u == '^' || u == 'v' {
				var ni, nj int
				switch u {
				case '<':
					ni, nj = i, j-1
				case '>':
					ni, nj = i, j+1
				case '^':
					ni, nj = i-1, j
				case 'v':
					ni, nj = i+1, j
				}
				if grid[ni][nj] == '#' {
					idx := i*m + j
					dist[idx] = 1
					queue = append(queue, idx)
				}
			}
		}
	}
	head := 0
	for head < len(queue) {
		v := queue[head]
		head++
		vi := v / m
		vj := v % m
		if vi > 0 {
			ui, uj := vi-1, vj
			if grid[ui][uj] == 'v' {
				u := ui*m + uj
				if dist[u] == 0 {
					dist[u] = dist[v] + 1
					queue = append(queue, u)
				}
			}
		}
		if vi+1 < n {
			ui, uj := vi+1, vj
			if grid[ui][uj] == '^' {
				u := ui*m + uj
				if dist[u] == 0 {
					dist[u] = dist[v] + 1
					queue = append(queue, u)
				}
			}
		}
		if vj > 0 {
			ui, uj := vi, vj-1
			if grid[ui][uj] == '>' {
				u := ui*m + uj
				if dist[u] == 0 {
					dist[u] = dist[v] + 1
					queue = append(queue, u)
				}
			}
		}
		if vj+1 < m {
			ui, uj := vi, vj+1
			if grid[ui][uj] == '<' {
				u := ui*m + uj
				if dist[u] == 0 {
					dist[u] = dist[v] + 1
					queue = append(queue, u)
				}
			}
		}
	}
	max1, max2 := 0, 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			c := grid[i][j]
			if c == '<' || c == '>' || c == '^' || c == 'v' {
				d := dist[i*m+j]
				if d == 0 {
					return -1
				}
				if d > max1 {
					max2 = max1
					max1 = d
				} else if d > max2 {
					max2 = d
				}
			}
		}
	}
	return max1 + max2
}

func genCase(r *rand.Rand) Test {
	n := r.Intn(4) + 3
	m := r.Intn(4) + 3
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if i == 0 || j == 0 || i == n-1 || j == m-1 {
				row[j] = '#'
			} else {
				opts := []byte{'<', '>', '^', 'v', '#'}
				row[j] = opts[r.Intn(len(opts))]
			}
		}
		grid[i] = row
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	out := fmt.Sprintf("%d", compute(grid))
	return Test{sb.String(), out}
}

func runCase(bin string, t Test) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(t.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(t.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 25; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
