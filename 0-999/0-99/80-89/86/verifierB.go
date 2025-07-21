package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type cell struct{ x, y int }

type shape []cell

var dirs = []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var validShapes map[string]struct{}

func canonical(s shape) string {
	best := ""
	trans := func(c cell, rot int, flip bool) cell {
		x, y := c.x, c.y
		// rotate
		for i := 0; i < rot; i++ {
			x, y = y, -x
		}
		if flip {
			x = -x
		}
		return cell{x, y}
	}
	for rot := 0; rot < 4; rot++ {
		for _, flip := range []bool{false, true} {
			tmp := make([]cell, len(s))
			for i, c := range s {
				tmp[i] = trans(c, rot, flip)
			}
			minX, minY := tmp[0].x, tmp[0].y
			for _, c := range tmp {
				if c.x < minX {
					minX = c.x
				}
				if c.y < minY {
					minY = c.y
				}
			}
			for i := range tmp {
				tmp[i].x -= minX
				tmp[i].y -= minY
			}
			sort.Slice(tmp, func(i, j int) bool {
				if tmp[i].x == tmp[j].x {
					return tmp[i].y < tmp[j].y
				}
				return tmp[i].x < tmp[j].x
			})
			var b strings.Builder
			for _, c := range tmp {
				b.WriteString(fmt.Sprintf("%d:%d,", c.x, c.y))
			}
			str := b.String()
			if best == "" || str < best {
				best = str
			}
		}
	}
	return best
}

func contains(s shape, c cell) bool {
	for _, x := range s {
		if x == c {
			return true
		}
	}
	return false
}

func generateShapes() {
	shapes := map[int]map[string]shape{}
	shapes[1] = map[string]shape{canonical(shape{{0, 0}}): {{0, 0}}}
	validShapes = make(map[string]struct{})
	for size := 1; size < 5; size++ {
		next := map[string]shape{}
		for _, sh := range shapes[size] {
			for _, c := range sh {
				for _, d := range dirs {
					nc := cell{c.x + d.x, c.y + d.y}
					if contains(sh, nc) {
						continue
					}
					ns := append([]cell(nil), sh...)
					ns = append(ns, nc)
					key := canonical(ns)
					if _, ok := next[key]; !ok {
						cp := append(shape(nil), ns...)
						next[key] = cp
						if len(cp) >= 2 {
							validShapes[key] = struct{}{}
						}
					}
				}
			}
		}
		shapes[size+1] = next
	}
}

func checkSolution(n, m int, orig []string, out []string) error {
	if len(out) != n {
		return fmt.Errorf("wrong number of lines")
	}
	for i := 0; i < n; i++ {
		if len(out[i]) != m {
			return fmt.Errorf("wrong line length")
		}
	}
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if orig[i][j] == '#' {
				if out[i][j] != '#' {
					return fmt.Errorf("cell %d %d should be #", i, j)
				}
				visited[i][j] = true
				continue
			}
			if out[i][j] < '0' || out[i][j] > '9' {
				return fmt.Errorf("invalid char at %d %d", i, j)
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			digit := out[i][j]
			q := []cell{{i, j}}
			visited[i][j] = true
			var cells shape
			for len(q) > 0 {
				c := q[0]
				q = q[1:]
				cells = append(cells, c)
				for _, d := range dirs {
					ni, nj := c.x+d.x, c.y+d.y
					if ni < 0 || ni >= n || nj < 0 || nj >= m {
						continue
					}
					if visited[ni][nj] {
						continue
					}
					if orig[ni][nj] == '#' {
						continue
					}
					if out[ni][nj] == digit {
						visited[ni][nj] = true
						q = append(q, cell{ni, nj})
					}
				}
			}
			if len(cells) < 2 || len(cells) > 5 {
				return fmt.Errorf("component size %d invalid", len(cells))
			}
			// convert to relative coords
			base := cells[0]
			rel := make(shape, len(cells))
			for idx, c := range cells {
				rel[idx] = cell{c.x - base.x, c.y - base.y}
			}
			if _, ok := validShapes[canonical(rel)]; !ok {
				return fmt.Errorf("invalid shape of size %d", len(cells))
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, bool) {
	if rng.Float64() < 0.1 {
		// unsolvable 1x1 board
		return "1 1\n.\n", false
	}
	n := rng.Intn(6) + 2
	m := rng.Intn(6) + 2
	board := make([]string, n)
	for i := 0; i < n; i++ {
		board[i] = strings.Repeat(".", m)
	}
	solvable := (n*m)%2 == 0
	return fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(board, "\n")), solvable
}

func runCase(bin, input string, solvable bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	if output == "-1" {
		if solvable {
			return fmt.Errorf("returned -1 on solvable case")
		}
		return nil
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	inputR := bufio.NewScanner(strings.NewReader(input))
	inputR.Scan() // first line n m
	n, m := 0, 0
	fmt.Sscan(inputR.Text(), &n, &m)
	orig := make([]string, n)
	for i := 0; i < n; i++ {
		inputR.Scan()
		orig[i] = inputR.Text()
	}
	return checkSolution(n, m, orig, lines)
}

func main() {
	generateShapes()
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, solv := generateCase(rng)
		if err := runCase(bin, in, solv); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
