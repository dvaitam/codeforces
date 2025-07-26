package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func components(grid [][]bool, n, m int) int {
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	comp := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] && !visited[i][j] {
				comp++
				queue := [][2]int{{i, j}}
				visited[i][j] = true
				for head := 0; head < len(queue); head++ {
					x, y := queue[head][0], queue[head][1]
					for _, d := range dirs {
						nx, ny := x+d[0], y+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] && !visited[nx][ny] {
							visited[nx][ny] = true
							queue = append(queue, [2]int{nx, ny})
						}
					}
				}
			}
		}
	}
	return comp
}

func hexRow(m int, s string) []bool {
	row := make([]bool, m)
	col := 0
	for i := 0; i < len(s); i++ {
		var v int
		ch := s[i]
		if ch >= '0' && ch <= '9' {
			v = int(ch - '0')
		} else {
			v = int(ch-'A') + 10
		}
		for b := 3; b >= 0; b-- {
			row[col] = (v>>uint(b))&1 == 1
			col++
		}
	}
	return row
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(10) + 1
		mWords := rand.Intn(3) + 1 // each word => 4 columns
		m := mWords * 4
		grid := make([][]bool, n)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			var rowStr strings.Builder
			for j := 0; j < mWords; j++ {
				val := rand.Intn(16)
				rowStr.WriteString(fmt.Sprintf("%X", val))
			}
			line := rowStr.String()
			fmt.Fprintln(&input, line)
			grid[i] = hexRow(m, line)
		}
		expected := components(grid, n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(expected) {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
