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

type DSU struct {
	parent []int
	sz     []int
	f      []bool
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), sz: make([]int, n), f: make([]bool, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *DSU) find(a int) int {
	if d.parent[a] == a {
		return a
	}
	p := d.parent[a]
	root := d.find(p)
	d.f[a] = d.f[a] != d.f[p]
	d.parent[a] = root
	return root
}

func (d *DSU) link(a, b int, parity bool) {
	if d.sz[a] < d.sz[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.sz[a] += d.sz[b]
	d.f[b] = parity
}

func (d *DSU) Unite(a, b int) {
	pa := d.find(a)
	pb := d.find(b)
	if pa != pb {
		parity := d.f[a] == d.f[b]
		d.link(pa, pb, parity)
	}
}

func solveCase(grid []string) (bool, [][]int) {
	n := len(grid)
	m := len(grid[0])
	dirs := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	id := func(i, j int) int { return i*m + j }
	dsu := NewDSU(n * m)
	res := make([][]int, n)
	for i := range res {
		res[i] = make([]int, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				continue
			}
			var ne []int
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					ne = append(ne, id(ni, nj))
				}
			}
			if len(ne)%2 == 1 {
				return false, nil
			}
			for k := 1; k < len(ne); k++ {
				dsu.Unite(ne[k-1], ne[k])
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '.' {
				continue
			}
			idx := id(i, j)
			dsu.find(idx)
			if dsu.f[idx] {
				res[i][j] = 1
			} else {
				res[i][j] = 4
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 'X' {
				continue
			}
			sum := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '.' {
					sum += res[ni][nj]
				}
			}
			res[i][j] = sum
		}
	}
	return true, res
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genGrid(rng *rand.Rand) []string {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}
	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			if rng.Intn(3) == 0 {
				grid[i][j] = 'X'
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func buildInput(grid []string) string {
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i] + "\n")
	}
	return sb.String()
}

func formatRes(ok bool, ans [][]int) string {
	if !ok {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i, row := range ans {
		for j, v := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		if i+1 < len(ans) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		grid := genGrid(rng)
		input := buildInput(grid)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		ok, ans := solveCase(grid)
		expected := formatRes(ok, ans)
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
