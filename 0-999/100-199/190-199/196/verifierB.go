package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func solveCase(n, m int, grid []string) string {
	NM := n * m
	visited := make([]bool, NM)
	offX := make([]int, NM)
	offY := make([]int, NM)
	start := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'S' {
				start = i*m + j
			}
		}
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queue := []int{start}
	visited[start] = true
	for q := 0; q < len(queue); q++ {
		v := queue[q]
		i := v / m
		j := v % m
		bx := offX[v]
		by := offY[v]
		for _, d := range dirs {
			ni := i + d[0]
			nj := j + d[1]
			nx, ny := bx, by
			if ni < 0 {
				ni += n
				nx--
			} else if ni >= n {
				ni -= n
				nx++
			}
			if nj < 0 {
				nj += m
				ny--
			} else if nj >= m {
				nj -= m
				ny++
			}
			if grid[ni][nj] == '#' {
				continue
			}
			u := ni*m + nj
			if !visited[u] {
				visited[u] = true
				offX[u] = nx
				offY[u] = ny
				queue = append(queue, u)
			} else if offX[u] != nx || offY[u] != ny {
				return "Yes"
			}
		}
	}
	return "No"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscanf(scan.Text(), "%d", &t)
	cases := make([]struct {
		n, m int
		grid []string
	}, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscanf(scan.Text(), "%d %d", &cases[i].n, &cases[i].m)
		cases[i].grid = make([]string, cases[i].n)
		for r := 0; r < cases[i].n; r++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			cases[i].grid[r] = scan.Text()
		}
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i, c := range cases {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		expected := solveCase(c.n, c.m, c.grid)
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
