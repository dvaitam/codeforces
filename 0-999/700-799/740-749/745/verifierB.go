package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pt struct{ x, y int }

func tryRect(H, W, h0, w0 int, grid0 []bool, rel []pt) bool {
	if H < h0 || W < w0 {
		return false
	}
	var dxs []int
	if H == h0 {
		dxs = []int{0}
	} else {
		dxs = []int{0, H - h0}
	}
	var dys []int
	if W == w0 {
		dys = []int{0}
	} else {
		dys = []int{0, W - w0}
	}
	for _, dx := range dxs {
		for _, dy := range dys {
			if dx == 0 && dy == 0 {
				continue
			}
			if dx >= h0 || dy >= w0 {
				return true
			}
			ok := true
			for _, p := range rel {
				i2 := p.x + dx
				j2 := p.y + dy
				if i2 >= 0 && i2 < h0 && j2 >= 0 && j2 < w0 {
					if grid0[i2*w0+j2] {
						ok = false
						break
					}
				}
			}
			if ok {
				return true
			}
		}
	}
	return false
}

func solve(n, m int, grid []string) string {
	pts := make([]pt, 0)
	minx, miny := n, m
	maxx, maxy := -1, -1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				pts = append(pts, pt{i, j})
				if i < minx {
					minx = i
				}
				if i > maxx {
					maxx = i
				}
				if j < miny {
					miny = j
				}
				if j > maxy {
					maxy = j
				}
			}
		}
	}
	h0 := maxx - minx + 1
	w0 := maxy - miny + 1
	k := len(pts)
	grid0 := make([]bool, h0*w0)
	rel := make([]pt, 0, k)
	for _, p := range pts {
		rx := p.x - minx
		ry := p.y - miny
		grid0[rx*w0+ry] = true
		rel = append(rel, pt{rx, ry})
	}
	t := 2 * k
	for H := 1; H*H <= t; H++ {
		if t%H != 0 {
			continue
		}
		W := t / H
		if tryRect(H, W, h0, w0, grid0, rel) {
			return "YES"
		}
		if H != W {
			if tryRect(W, H, h0, w0, grid0, rel) {
				return "YES"
			}
		}
	}
	return "NO"
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type testCase struct {
	input    string
	expected string
}

var dx = []int{-1, 1, 0, 0}
var dy = []int{0, 0, -1, 1}

func randomPiece(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	sx := rng.Intn(n)
	sy := rng.Intn(m)
	grid[sx][sy] = 'X'
	cells := []pt{{sx, sy}}
	total := rng.Intn(n*m) + 1
	for len(cells) < total {
		p := cells[rng.Intn(len(cells))]
		dir := rng.Intn(4)
		nx := p.x + dx[dir]
		ny := p.y + dy[dir]
		if nx < 0 || nx >= n || ny < 0 || ny >= m || grid[nx][ny] == 'X' {
			continue
		}
		grid[nx][ny] = 'X'
		cells = append(cells, pt{nx, ny})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	strGrid := make([]string, n)
	for i := 0; i < n; i++ {
		strGrid[i] = string(grid[i])
		sb.WriteString(strGrid[i] + "\n")
	}
	exp := solve(n, m, strGrid)
	return testCase{input: sb.String(), expected: exp}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(746))

	cases := []testCase{
		{input: "1 1\nX\n", expected: solve(1, 1, []string{"X"})},
		{input: "2 1\nX\nX\n", expected: solve(2, 1, []string{"X", "X"})},
		{input: "1 2\nXX\n", expected: solve(1, 2, []string{"XX"})},
		{input: "2 2\nXX\nXX\n", expected: solve(2, 2, []string{"XX", "XX"})},
		{input: "2 2\nXX\nX.\n", expected: solve(2, 2, []string{"XX", "X."})},
	}
	for len(cases) < 100 {
		cases = append(cases, randomPiece(rng))
	}

	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, tc.expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
