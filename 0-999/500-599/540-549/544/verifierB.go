package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countIslands(grid []string) int {
	n := len(grid)
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	dirs := []struct{ x, y int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	q := make([][2]int, 0)
	islands := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != 'L' || visited[i][j] {
				continue
			}
			islands++
			visited[i][j] = true
			q = q[:0]
			q = append(q, [2]int{i, j})
			for len(q) > 0 {
				cur := q[0]
				q = q[1:]
				for _, d := range dirs {
					ni, nj := cur[0]+d.x, cur[1]+d.y
					if ni >= 0 && ni < n && nj >= 0 && nj < n && !visited[ni][nj] && grid[ni][nj] == 'L' {
						visited[ni][nj] = true
						q = append(q, [2]int{ni, nj})
					}
				}
			}
		}
	}
	return islands
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d malformed line\n", idx)
			os.Exit(1)
		}
		var n, k int
		fmt.Sscan(parts[0], &n)
		fmt.Sscan(parts[1], &k)
		input := fmt.Sprintf("%d %d\n", n, k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Printf("test %d: empty output\n", idx)
			os.Exit(1)
		}
		ans := strings.TrimSpace(lines[0])
		maxIslands := (n*n + 1) / 2
		if ans == "NO" {
			if k <= maxIslands {
				fmt.Printf("test %d: answer should be YES\n", idx)
				os.Exit(1)
			}
			continue
		}
		if ans != "YES" {
			fmt.Printf("test %d: first line must be YES or NO\n", idx)
			os.Exit(1)
		}
		if k > maxIslands {
			fmt.Printf("test %d: answer should be NO\n", idx)
			os.Exit(1)
		}
		if len(lines)-1 != n {
			fmt.Printf("test %d: expected %d lines, got %d\n", idx, n, len(lines)-1)
			os.Exit(1)
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if len(lines[i+1]) != n {
				fmt.Printf("test %d: line %d wrong length\n", idx, i+1)
				os.Exit(1)
			}
			for _, ch := range lines[i+1] {
				if ch != 'L' && ch != 'S' {
					fmt.Printf("test %d: invalid character\n", idx)
					os.Exit(1)
				}
			}
			grid[i] = lines[i+1]
		}
		islands := countIslands(grid)
		if islands != k {
			fmt.Printf("test %d: expected %d islands got %d\n", idx, k, islands)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
