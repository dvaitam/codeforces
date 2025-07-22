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

var n, m int
var grid [][]byte

func inBounds(x, y int) bool {
	return x >= 0 && x < n && y >= 0 && y < m
}

func applyMove(mask uint32, x, y int) uint32 {
	idx := x*m + y
	bit := uint32(1) << idx
	if mask&bit == 0 {
		return mask
	}
	mask &^= bit
	var dirs [][2]int
	switch grid[x][y] {
	case 'L':
		dirs = [][2]int{{1, -1}, {-1, 1}}
	case 'R':
		dirs = [][2]int{{-1, -1}, {1, 1}}
	case 'X':
		dirs = [][2]int{{1, -1}, {-1, 1}, {-1, -1}, {1, 1}}
	}
	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]
		for inBounds(nx, ny) {
			idx := nx*m + ny
			b := uint32(1) << idx
			if mask&b == 0 {
				break
			}
			mask &^= b
			nx += d[0]
			ny += d[1]
		}
	}
	return mask
}

var memo map[uint32]int

func grundy(mask uint32) int {
	if mask == 0 {
		return 0
	}
	if v, ok := memo[mask]; ok {
		return v
	}
	moves := make(map[int]struct{})
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			idx := i*m + j
			if mask&(1<<idx) != 0 {
				nm := applyMove(mask, i, j)
				g := grundy(nm)
				moves[g] = struct{}{}
			}
		}
	}
	mex := 0
	for {
		if _, ok := moves[mex]; !ok {
			break
		}
		mex++
	}
	memo[mask] = mex
	return mex
}

func expectedWinner(g [][]byte) string {
	n = len(g)
	m = len(g[0])
	grid = g
	memo = make(map[uint32]int)
	mask := uint32(0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			mask |= 1 << (i*m + j)
		}
	}
	if grundy(mask) != 0 {
		return "WIN"
	}
	return "LOSE"
}

func generateCase(rng *rand.Rand) (string, string) {
	n = rng.Intn(3) + 1
	m = rng.Intn(3) + 1
	grid = make([][]byte, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	chars := []byte{'L', 'R', 'X'}
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = chars[rng.Intn(len(chars))]
		}
		grid[i] = row
		sb.Write(row)
		sb.WriteByte('\n')
	}
	exp := expectedWinner(grid)
	return sb.String(), exp
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
