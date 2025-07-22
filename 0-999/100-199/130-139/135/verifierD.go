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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(grid [][]byte) int {
	n := len(grid)
	if n == 0 {
		return 0
	}
	m := len(grid[0])
	outside := make([][]bool, n)
	for i := range outside {
		outside[i] = make([]bool, m)
	}
	dq := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for _, j := range []int{0, m - 1} {
			if grid[i][j] == '0' && !outside[i][j] {
				outside[i][j] = true
				dq = append(dq, [2]int{i, j})
			}
		}
	}
	for j := 0; j < m; j++ {
		for _, i := range []int{0, n - 1} {
			if grid[i][j] == '0' && !outside[i][j] {
				outside[i][j] = true
				dq = append(dq, [2]int{i, j})
			}
		}
	}
	for idx := 0; idx < len(dq); idx++ {
		r, c := dq[idx][0], dq[idx][1]
		for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '0' && !outside[nr][nc] {
				outside[nr][nc] = true
				dq = append(dq, [2]int{nr, nc})
			}
		}
	}
	dr := [4]int{-1, 0, 1, 0}
	dc := [4]int{0, 1, 0, -1}
	visited := make([][][]bool, n)
	for i := range visited {
		visited[i] = make([][]bool, m)
		for j := range visited[i] {
			visited[i][j] = make([]bool, 4)
		}
	}
	best := 0
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			if grid[r][c] != '1' {
				continue
			}
			for d := 0; d < 4; d++ {
				nr, nc := r+dr[d], c+dc[d]
				if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] != '1' {
					continue
				}
				if visited[r][c][d] {
					continue
				}
				cr, cc, dir := r, c, d
				startR, startC, startDir := r, c, d
				length := 0
				infinite := false
				for {
					visited[cr][cc][dir] = true
					nr, nc := cr+dr[dir], cc+dc[dir]
					if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] != '1' {
						infinite = true
						break
					}
					found := false
					for _, turn := range []int{1, 0, 3, 2} {
						nd := (dir + turn) & 3
						wr, wc := nr+dr[nd], nc+dc[nd]
						if wr >= 0 && wr < n && wc >= 0 && wc < m && grid[wr][wc] == '1' {
							cr, cc, dir = nr, nc, nd
							found = true
							break
						}
					}
					if !found {
						infinite = true
						break
					}
					length++
					if cr == startR && cc == startC && dir == startDir {
						break
					}
				}
				if !infinite && length > 0 && length > best {
					best = length
				}
			}
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	grid := make([][]byte, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		grid[i] = row
		sb.WriteString(string(row) + "\n")
	}
	expect := solveCase(grid)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(out, &val); err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, val, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
