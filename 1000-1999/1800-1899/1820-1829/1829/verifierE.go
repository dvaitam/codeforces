package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pair struct{ x, y int }

func solveE(grid [][]int) string {
	n := len(grid)
	m := len(grid[0])
	dx := []int{1, -1, 0, 0}
	dy := []int{0, 0, 1, -1}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	best := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 || visited[i][j] {
				continue
			}
			sum := 0
			q := []pair{{i, j}}
			visited[i][j] = true
			for head := 0; head < len(q); head++ {
				p := q[head]
				sum += grid[p.x][p.y]
				for dir := 0; dir < 4; dir++ {
					nx := p.x + dx[dir]
					ny := p.y + dy[dir]
					if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] && grid[nx][ny] > 0 {
						visited[nx][ny] = true
						q = append(q, pair{nx, ny})
					}
				}
			}
			if sum > best {
				best = sum
			}
		}
	}
	return fmt.Sprint(best)
}

func genTestsE() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for tc := 0; tc < t; tc++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		fmt.Fprintf(&input, "%d %d\n", n, m)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				grid[i][j] = rand.Intn(11)
				if j+1 == m {
					fmt.Fprintln(&input, grid[i][j])
				} else {
					fmt.Fprint(&input, grid[i][j], " ")
				}
			}
		}
		expected[tc] = solveE(grid)
	}
	return expected, input.String()
}

func runBinary(path, in string) ([]string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(&out)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsE()
	lines, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if len(lines) != len(expected) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expected), len(lines))
		os.Exit(1)
	}
	for i, exp := range expected {
		if lines[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
