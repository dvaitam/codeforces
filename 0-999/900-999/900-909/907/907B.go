package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	rows := make([]string, 0, 9)
	for in.Scan() {
		line := in.Text()
		if len(line) == 0 {
			continue
		}
		clean := strings.ReplaceAll(line, " ", "")
		rows = append(rows, clean)
		if len(rows) == 9 {
			break
		}
	}
	if len(rows) != 9 {
		return
	}

	var x, y int
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		fmt.Sscanf(line, "%d %d", &x, &y)
		break
	}

	grid := make([][]byte, 9)
	for i := 0; i < 9; i++ {
		grid[i] = []byte(rows[i])
	}

	targetRow := ((x - 1) % 3) * 3
	targetCol := ((y - 1) % 3) * 3
	hasEmpty := false
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[targetRow+i][targetCol+j] == '.' {
				hasEmpty = true
			}
		}
	}

	if hasEmpty {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if grid[targetRow+i][targetCol+j] == '.' {
					grid[targetRow+i][targetCol+j] = '!'
				}
			}
		}
	} else {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if grid[i][j] == '.' {
					grid[i][j] = '!'
				}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			out.WriteByte(grid[i][j])
			if j%3 == 2 && j != 8 {
				out.WriteByte(' ')
			}
		}
		out.WriteByte('\n')
		if i%3 == 2 && i != 8 {
			out.WriteByte('\n')
		}
	}
	out.Flush()
}
