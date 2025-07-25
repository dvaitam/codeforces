package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type pair struct{ x, y int }

func expectedD(grid [][]byte, queries []pair) []int {
	n, m := len(grid), len(grid[0])
	comp := make([][]int, n)
	for i := range comp {
		comp[i] = make([]int, m)
	}
	pictures := []int{0}
	compID := 0
	dirs := []int{-1, 0, 1, 0, -1}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '.' || comp[i][j] != 0 {
				continue
			}
			compID++
			queue := []pair{{i, j}}
			comp[i][j] = compID
			pics := 0
			for head := 0; head < len(queue); head++ {
				cur := queue[head]
				x, y := cur.x, cur.y
				for d := 0; d < 4; d++ {
					nx := x + dirs[d]
					ny := y + dirs[d+1]
					if nx < 0 || nx >= n || ny < 0 || ny >= m {
						continue
					}
					if grid[nx][ny] == '*' {
						pics++
					} else if comp[nx][ny] == 0 {
						comp[nx][ny] = compID
						queue = append(queue, pair{nx, ny})
					}
				}
			}
			pictures = append(pictures, pics)
		}
	}

	res := make([]int, len(queries))
	for i, q := range queries {
		res[i] = pictures[comp[q.x][q.y]]
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	m := rng.Intn(6) + 3
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if i == 0 || j == 0 || i == n-1 || j == m-1 {
				row[j] = '*'
			} else if rng.Intn(4) == 0 {
				row[j] = '*'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}
	k := rng.Intn(min(n*m, 5)) + 1
	queries := make([]pair, k)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]) + "\n")
	}
	empties := make([]pair, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				empties = append(empties, pair{i, j})
			}
		}
	}
	if len(empties) == 0 {
		empties = append(empties, pair{1, 1})
		grid[1][1] = '.'
	}
	for i := 0; i < k; i++ {
		p := empties[rng.Intn(len(empties))]
		queries[i] = p
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x+1, p.y+1))
	}
	ans := expectedD(grid, queries)
	var out strings.Builder
	for i, v := range ans {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(strconv.Itoa(v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected\n%s\ngot\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
