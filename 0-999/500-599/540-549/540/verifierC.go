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

type caseC struct {
	input    string
	expected string
}

func solveCase(n, m int, grid [][]byte, r1, c1, r2, c2 int) string {
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	type pt struct{ x, y int }
	q := []pt{{r1, c1}}
	visited[r1][c1] = true
	dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for h := 0; h < len(q); h++ {
		p := q[h]
		for _, d := range dirs {
			nx, ny := p.x+d.x, p.y+d.y
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if visited[nx][ny] {
				continue
			}
			if grid[nx][ny] == '.' || (nx == r2 && ny == c2) {
				visited[nx][ny] = true
				q = append(q, pt{nx, ny})
			}
		}
	}
	if !visited[r2][c2] {
		return "NO"
	}
	cnt := 0
	dirs = []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, d := range dirs {
		nx, ny := r2+d.x, c2+d.y
		if nx < 0 || nx >= n || ny < 0 || ny >= m {
			continue
		}
		if grid[nx][ny] == '.' || (nx == r1 && ny == c1) {
			cnt++
		}
	}
	if r1 == r2 && c1 == c2 {
		if cnt >= 1 {
			return "YES"
		}
		return "NO"
	}
	if grid[r2][c2] == 'X' {
		return "YES"
	}
	if cnt >= 2 {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) caseC {
	n := rng.Intn(9) + 1
	m := rng.Intn(9) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = 'X'
			}
		}
		grid[i] = row
	}
	r1 := rng.Intn(n)
	c1 := rng.Intn(m)
	grid[r1][c1] = 'X'
	r2 := rng.Intn(n)
	c2 := rng.Intn(m)
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		b.WriteString(string(grid[i]))
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "%d %d\n%d %d\n", r1+1, c1+1, r2+1, c2+1)
	expected := solveCase(n, m, grid, r1, c1, r2, c2)
	return caseC{b.String(), expected}
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		got, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != c.expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, c.expected, got, c.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
