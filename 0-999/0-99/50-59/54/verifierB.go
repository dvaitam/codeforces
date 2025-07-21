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

type pair struct{ x, y int }

func repr(piece [][]byte) string {
	var b strings.Builder
	for i, row := range piece {
		b.Write(row)
		if i+1 < len(piece) {
			b.WriteByte('|')
		}
	}
	return b.String()
}

func rot180(p [][]byte) [][]byte {
	x := len(p)
	y := len(p[0])
	r := make([][]byte, x)
	for i := 0; i < x; i++ {
		r[i] = make([]byte, y)
		for j := 0; j < y; j++ {
			r[i][j] = p[x-1-i][y-1-j]
		}
	}
	return r
}

func rot90(p [][]byte) [][]byte {
	x := len(p)
	y := len(p[0])
	r := make([][]byte, x)
	for i := range r {
		r[i] = make([]byte, y)
	}
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r[j][x-1-i] = p[i][j]
		}
	}
	return r
}

func rot270(p [][]byte) [][]byte { return rot180(rot90(p)) }

func good(grid [][]byte, X, Y int) bool {
	A := len(grid)
	B := len(grid[0])
	rows := A / X
	cols := B / Y
	seen := make(map[string]struct{})
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			piece := make([][]byte, X)
			for k := 0; k < X; k++ {
				piece[k] = grid[i*X+k][j*Y : j*Y+Y]
			}
			reps := []string{repr(piece), repr(rot180(piece))}
			if X == Y {
				reps = append(reps, repr(rot90(piece)), repr(rot270(piece)))
			}
			m := reps[0]
			for _, s := range reps[1:] {
				if s < m {
					m = s
				}
			}
			if _, ok := seen[m]; ok {
				return false
			}
			seen[m] = struct{}{}
		}
	}
	return true
}

func solve(grid [][]byte) (int, pair) {
	A := len(grid)
	B := len(grid[0])
	var divA, divB []int
	for x := 1; x <= A; x++ {
		if A%x == 0 {
			divA = append(divA, x)
		}
	}
	for y := 1; y <= B; y++ {
		if B%y == 0 {
			divB = append(divB, y)
		}
	}
	var res []pair
	for _, X := range divA {
		for _, Y := range divB {
			if good(grid, X, Y) {
				res = append(res, pair{X, Y})
			}
		}
	}
	best := res[0]
	bestArea := best.x * best.y
	for _, p := range res[1:] {
		area := p.x * p.y
		if area < bestArea || (area == bestArea && p.x < best.x) {
			best = p
			bestArea = area
		}
	}
	return len(res), best
}

func runCase(exe string, input string, expCount int, expPair pair) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var c, x, y int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &c, &x, &y); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if c != expCount || x != expPair.x || y != expPair.y {
		return fmt.Errorf("expected %d %d %d got %s", expCount, expPair.x, expPair.y, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 100; i++ {
		A := rng.Intn(5) + 1
		B := rng.Intn(5) + 1
		grid := make([][]byte, A)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", A, B))
		for r := 0; r < A; r++ {
			row := make([]byte, B)
			for c := 0; c < B; c++ {
				row[c] = letters[rng.Intn(26)]
			}
			grid[r] = row
			sb.Write(row)
			sb.WriteByte('\n')
		}
		input := sb.String()
		cnt, best := solve(grid)
		if err := runCase(exe, input, cnt, best); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
