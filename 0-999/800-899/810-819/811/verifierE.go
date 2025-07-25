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

type point struct{ x, y int }

func countComponents(grid [][]int, n, m int, l, r int) int {
	w := r - l + 1
	vis := make([][]bool, n)
	for i := 0; i < n; i++ {
		vis[i] = make([]bool, w)
	}
	comp := 0
	dirs := []struct{ dx, dy int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < w; j++ {
			if vis[i][j] {
				continue
			}
			comp++
			val := grid[i][j+l]
			q := []point{{i, j}}
			vis[i][j] = true
			for len(q) > 0 {
				p := q[0]
				q = q[1:]
				for _, d := range dirs {
					nx, ny := p.x+d.dx, p.y+d.dy
					if nx < 0 || nx >= n || ny < 0 || ny >= w {
						continue
					}
					if vis[nx][ny] {
						continue
					}
					if grid[nx][ny+l] != val {
						continue
					}
					vis[nx][ny] = true
					q = append(q, point{nx, ny})
				}
			}
		}
	}
	return comp
}

func solveCase(n, m, q int, grid [][]int, queries [][2]int) []int {
	ans := make([]int, q)
	for i, qu := range queries {
		ans[i] = countComponents(grid, n, m, qu[0]-1, qu[1]-1)
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(6) + 1
	q := rng.Intn(5) + 1
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = rng.Intn(3) + 1
		}
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := rng.Intn(m-l+1) + l
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", grid[i][j])
		}
		sb.WriteByte('\n')
	}
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qu[0], qu[1])
	}
	answers := solveCase(n, m, q, grid, queries)
	var exp strings.Builder
	for i, a := range answers {
		if i > 0 {
			exp.WriteByte('\n')
		}
		fmt.Fprintf(&exp, "%d", a)
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
