package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	n    int
	m    int
	k    int
	grid [][]int
}

func genTestsB() []testCaseB {
	rand.Seed(2)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(3) + 2 // 2..4
		m := rand.Intn(3) + 2
		k := rand.Intn(4)
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			grid[r] = make([]int, m)
			for c := 0; c < m; c++ {
				grid[r][c] = rand.Intn(2)
			}
		}
		tests[i] = testCaseB{n, m, k, grid}
	}
	return tests
}

type cell struct{ r, c int }

func valid(grid [][]int) bool {
	n := len(grid)
	m := len(grid[0])
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, m)
	}
	dirs := []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if vis[i][j] {
				continue
			}
			v := grid[i][j]
			q := []cell{{i, j}}
			vis[i][j] = true
			cells := []cell{{i, j}}
			minR, maxR := i, i
			minC, maxC := j, j
			for len(q) > 0 {
				u := q[0]
				q = q[1:]
				if u.r < minR {
					minR = u.r
				}
				if u.r > maxR {
					maxR = u.r
				}
				if u.c < minC {
					minC = u.c
				}
				if u.c > maxC {
					maxC = u.c
				}
				for _, d := range dirs {
					nr, nc := u.r+d.r, u.c+d.c
					if nr >= 0 && nr < n && nc >= 0 && nc < m && !vis[nr][nc] && grid[nr][nc] == v {
						vis[nr][nc] = true
						q = append(q, cell{nr, nc})
						cells = append(cells, cell{nr, nc})
					}
				}
			}
			area := (maxR - minR + 1) * (maxC - minC + 1)
			if area != len(cells) {
				return false
			}
			for rr := minR; rr <= maxR; rr++ {
				for cc := minC; cc <= maxC; cc++ {
					if grid[rr][cc] != v {
						return false
					}
				}
			}
		}
	}
	return true
}

func minChanges(grid [][]int, k int) int {
	n := len(grid)
	m := len(grid[0])
	cells := make([]cell, 0, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cells = append(cells, cell{i, j})
		}
	}
	best := k + 1
	N := len(cells)
	for mask := 0; mask < 1<<N; mask++ {
		cnt := bits.OnesCount(uint(mask))
		if cnt > k {
			continue
		}
		changed := make([][]int, n)
		for i := 0; i < n; i++ {
			changed[i] = make([]int, m)
			copy(changed[i], grid[i])
		}
		for b := 0; b < N; b++ {
			if mask&(1<<b) != 0 {
				c := cells[b]
				changed[c.r][c.c] ^= 1
			}
		}
		if valid(changed) {
			if cnt < best {
				best = cnt
			}
		}
	}
	if best > k {
		return -1
	}
	return best
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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, tc.k)
		for r := 0; r < tc.n; r++ {
			for c := 0; c < tc.m; c++ {
				if c > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, tc.grid[r][c])
			}
			input.WriteByte('\n')
		}
		expect := minChanges(tc.grid, tc.k)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
