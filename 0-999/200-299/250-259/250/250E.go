package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
	}

	row, col := 0, 0
	dir := 1
	seconds := int64(0)
	token := 1
	visited := make([]int, n*m*2)

	for {
		if row == n-1 {
			fmt.Println(seconds)
			return
		}

		dirIdx := 0
		if dir == 1 {
			dirIdx = 1
		}
		idx := ((row*m)+col)*2 + dirIdx
		if visited[idx] == token {
			fmt.Println("Never")
			return
		}
		visited[idx] = token

		if row+1 < n && grid[row+1][col] == '.' {
			row++
			seconds++
			continue
		}

		nextCol := col + dir
		seconds++
		if nextCol < 0 || nextCol >= m {
			dir = -dir
			continue
		}

		cell := grid[row][nextCol]
		if cell == '.' {
			col = nextCol
			continue
		}
		if cell == '+' {
			grid[row][nextCol] = '.'
			dir = -dir
			token++
			continue
		}
		dir = -dir
	}
}

