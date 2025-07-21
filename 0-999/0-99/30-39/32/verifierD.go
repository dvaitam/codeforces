package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Cross struct {
	r, i, j int
}

func expected(n, m int, k int64, grid []string) string {
	up := make([][]int, n)
	down := make([][]int, n)
	left := make([][]int, n)
	right := make([][]int, n)
	for i := 0; i < n; i++ {
		up[i] = make([]int, m)
		down[i] = make([]int, m)
		left[i] = make([]int, m)
		right[i] = make([]int, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' {
				if i > 0 {
					up[i][j] = up[i-1][j] + 1
				} else {
					up[i][j] = 1
				}
				if j > 0 {
					left[i][j] = left[i][j-1] + 1
				} else {
					left[i][j] = 1
				}
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if grid[i][j] == '*' {
				if i < n-1 {
					down[i][j] = down[i+1][j] + 1
				} else {
					down[i][j] = 1
				}
				if j < m-1 {
					right[i][j] = right[i][j+1] + 1
				} else {
					right[i][j] = 1
				}
			}
		}
	}
	var crosses []Cross
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' {
				rmax := up[i][j] - 1
				if down[i][j]-1 < rmax {
					rmax = down[i][j] - 1
				}
				if left[i][j]-1 < rmax {
					rmax = left[i][j] - 1
				}
				if right[i][j]-1 < rmax {
					rmax = right[i][j] - 1
				}
				for r := 1; r <= rmax; r++ {
					crosses = append(crosses, Cross{r, i + 1, j + 1})
				}
			}
		}
	}
	sort.Slice(crosses, func(a, b int) bool {
		if crosses[a].r != crosses[b].r {
			return crosses[a].r < crosses[b].r
		}
		if crosses[a].i != crosses[b].i {
			return crosses[a].i < crosses[b].i
		}
		return crosses[a].j < crosses[b].j
	})
	if int64(len(crosses)) < k {
		return "-1"
	}
	c := crosses[k-1]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", c.i, c.j))
	sb.WriteString(fmt.Sprintf("%d %d\n", c.i-c.r, c.j))
	sb.WriteString(fmt.Sprintf("%d %d\n", c.i+c.r, c.j))
	sb.WriteString(fmt.Sprintf("%d %d\n", c.i, c.j-c.r))
	sb.WriteString(fmt.Sprintf("%d %d", c.i, c.j+c.r))
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 3
	m := rng.Intn(7) + 3
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		var row strings.Builder
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row.WriteByte('.')
			} else {
				row.WriteByte('*')
			}
		}
		grid[i] = row.String()
	}
	crosses := make([]Cross, 0)
	up := make([][]int, n)
	down := make([][]int, n)
	left := make([][]int, n)
	right := make([][]int, n)
	for i := 0; i < n; i++ {
		up[i] = make([]int, m)
		down[i] = make([]int, m)
		left[i] = make([]int, m)
		right[i] = make([]int, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' {
				if i > 0 {
					up[i][j] = up[i-1][j] + 1
				} else {
					up[i][j] = 1
				}
				if j > 0 {
					left[i][j] = left[i][j-1] + 1
				} else {
					left[i][j] = 1
				}
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if grid[i][j] == '*' {
				if i < n-1 {
					down[i][j] = down[i+1][j] + 1
				} else {
					down[i][j] = 1
				}
				if j < m-1 {
					right[i][j] = right[i][j+1] + 1
				} else {
					right[i][j] = 1
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' {
				rmax := up[i][j] - 1
				if down[i][j]-1 < rmax {
					rmax = down[i][j] - 1
				}
				if left[i][j]-1 < rmax {
					rmax = left[i][j] - 1
				}
				if right[i][j]-1 < rmax {
					rmax = right[i][j] - 1
				}
				for r := 1; r <= rmax; r++ {
					crosses = append(crosses, Cross{r, i + 1, j + 1})
				}
			}
		}
	}
	k := int64(rng.Intn(len(crosses)+3) + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	exp := expected(n, m, k, grid)
	return sb.String(), exp
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
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
