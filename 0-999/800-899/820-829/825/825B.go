package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkWin(b [][]byte) bool {
	dirs := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if b[i][j] != 'X' {
				continue
			}
			for _, d := range dirs {
				dx, dy := d[0], d[1]
				ok := true
				for k := 1; k < 5; k++ {
					x := i + dx*k
					y := j + dy*k
					if x < 0 || x >= 10 || y < 0 || y >= 10 || b[x][y] != 'X' {
						ok = false
						break
					}
				}
				if ok {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	board := make([][]byte, 10)
	for i := 0; i < 10; i++ {
		var line string
		fmt.Fscan(reader, &line)
		board[i] = []byte(line)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if board[i][j] == '.' {
				board[i][j] = 'X'
				if checkWin(board) {
					fmt.Println("YES")
					return
				}
				board[i][j] = '.'
			}
		}
	}
	fmt.Println("NO")
}
