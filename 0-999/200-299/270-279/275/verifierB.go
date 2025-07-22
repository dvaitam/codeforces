package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isConvex(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	rowPS := make([][]int, n)
	for i := 0; i < n; i++ {
		rowPS[i] = make([]int, m+1)
		for j := 0; j < m; j++ {
			rowPS[i][j+1] = rowPS[i][j]
			if grid[i][j] == 'B' {
				rowPS[i][j+1]++
			}
		}
	}
	colPS := make([][]int, m)
	for j := 0; j < m; j++ {
		colPS[j] = make([]int, n+1)
		for i := 0; i < n; i++ {
			colPS[j][i+1] = colPS[j][i]
			if grid[i][j] == 'B' {
				colPS[j][i+1]++
			}
		}
	}
	type cell struct{ r, c int }
	blacks := make([]cell, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'B' {
				blacks = append(blacks, cell{i, j})
			}
		}
	}
	k := len(blacks)
	for a := 0; a < k; a++ {
		for b := a + 1; b < k; b++ {
			r1, c1 := blacks[a].r, blacks[a].c
			r2, c2 := blacks[b].r, blacks[b].c
			ok := false
			if r1 == r2 {
				cmin, cmax := c1, c2
				if cmin > cmax {
					cmin, cmax = cmax, cmin
				}
				if rowPS[r1][cmax+1]-rowPS[r1][cmin] == cmax-cmin+1 {
					ok = true
				}
			}
			if !ok && c1 == c2 {
				rmin, rmax := r1, r2
				if rmin > rmax {
					rmin, rmax = rmax, rmin
				}
				if colPS[c1][rmax+1]-colPS[c1][rmin] == rmax-rmin+1 {
					ok = true
				}
			}
			if !ok {
				cmin, cmax := c1, c2
				if cmin > cmax {
					cmin, cmax = cmax, cmin
				}
				rmin, rmax := r1, r2
				if rmin > rmax {
					rmin, rmax = rmax, rmin
				}
				seg1 := rowPS[r1][cmax+1]-rowPS[r1][cmin] == cmax-cmin+1
				seg2 := colPS[c2][rmax+1]-colPS[c2][rmin] == rmax-rmin+1
				if seg1 && seg2 {
					ok = true
				}
			}
			if !ok {
				cmin, cmax := c1, c2
				if cmin > cmax {
					cmin, cmax = cmax, cmin
				}
				rmin, rmax := r1, r2
				if rmin > rmax {
					rmin, rmax = rmax, rmin
				}
				seg1 := colPS[c1][rmax+1]-colPS[c1][rmin] == rmax-rmin+1
				seg2 := rowPS[r2][cmax+1]-rowPS[r2][cmin] == cmax-cmin+1
				if seg1 && seg2 {
					ok = true
				}
			}
			if !ok {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([][]byte, n)
		bcount := 0
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if rng.Intn(2) == 0 {
					row[c] = 'W'
				} else {
					row[c] = 'B'
					bcount++
				}
			}
			grid[r] = row
		}
		if bcount == 0 {
			grid[0][0] = 'B'
			bcount = 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for r := 0; r < n; r++ {
			sb.WriteString(string(grid[r]))
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := "NO"
		if isConvex(grid) {
			expected = "YES"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
