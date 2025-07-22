package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveB(r, c int, grid []string) string {
	var si, sj, ei, ej int
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ch := grid[i][j]
			if ch == 'S' {
				si, sj = i, j
			} else if ch == 'E' {
				ei, ej = i, j
			}
		}
	}
	n := r * c
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	idx0 := ei*c + ej
	dist[idx0] = 0
	q := list.New()
	q.PushBack(idx0)
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for q.Len() > 0 {
		v := q.Remove(q.Front()).(int)
		x := v / c
		y := v % c
		d := dist[v]
		for _, dir := range dirs {
			ni := x + dir[0]
			nj := y + dir[1]
			if ni < 0 || ni >= r || nj < 0 || nj >= c {
				continue
			}
			ch := grid[ni][nj]
			if ch == 'T' {
				continue
			}
			idx := ni*c + nj
			if dist[idx] == -1 {
				dist[idx] = d + 1
				q.PushBack(idx)
			}
		}
	}
	startDist := dist[si*c+sj]
	total := 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ch := grid[i][j]
			if ch >= '0' && ch <= '9' {
				d := dist[i*c+j]
				if d != -1 && d <= startDist {
					total += int(ch - '0')
				}
			}
		}
	}
	return fmt.Sprintf("%d", total)
}

func generateCaseB(rng *rand.Rand) (string, string) {
	r := rng.Intn(5) + 1
	c := rng.Intn(5) + 1
	grid := make([]string, r)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	si := rng.Intn(r)
	sj := rng.Intn(c)
	ei := rng.Intn(r)
	ej := rng.Intn(c)
	for ei == si && ej == sj {
		ei = rng.Intn(r)
		ej = rng.Intn(c)
	}
	for i := 0; i < r; i++ {
		row := make([]byte, c)
		for j := 0; j < c; j++ {
			if i == si && j == sj {
				row[j] = 'S'
			} else if i == ei && j == ej {
				row[j] = 'E'
			} else {
				if rng.Intn(4) == 0 {
					row[j] = byte('0' + rng.Intn(10))
				} else {
					row[j] = '.'
				}
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	out := solveB(r, c, grid)
	return sb.String(), out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
