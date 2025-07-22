package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	maze := make([]string, n)
	var startX, startY int
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &maze[i])
		if j := indexByte(maze[i], 'S'); j != -1 {
			startX, startY = i, j
		}
	}
	var instr string
	fmt.Fscan(reader, &instr)

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	count := 0
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			if b == a {
				continue
			}
			for c := 0; c < 4; c++ {
				if c == a || c == b {
					continue
				}
				for d := 0; d < 4; d++ {
					if d == a || d == b || d == c {
						continue
					}
					mapping := [4]int{a, b, c, d}
					x, y := startX, startY
					success := false
					for i := 0; i < len(instr); i++ {
						dIdx := mapping[int(instr[i]-'0')]
						nx := x + dirs[dIdx][0]
						ny := y + dirs[dIdx][1]
						if nx < 0 || nx >= n || ny < 0 || ny >= m {
							success = false
							break
						}
						cell := maze[nx][ny]
						if cell == '#' {
							success = false
							break
						}
						if cell == 'E' {
							success = true
							break
						}
						x, y = nx, ny
					}
					if success {
						count++
					}
				}
			}
		}
	}
	fmt.Println(count)
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
