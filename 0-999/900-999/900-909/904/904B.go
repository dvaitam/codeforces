package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	board := make([][]byte, 0, 9)
	for len(board) < 9 {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			break
		}
		line = strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		line = strings.ReplaceAll(line, " ", "")
		if len(line) == 0 {
			continue
		}
		row := []byte(line)
		board = append(board, row)
	}
	var x, y int
	fmt.Fscan(reader, &x, &y)
	// convert to 0-index
	targetRow := (x - 1) % 3
	targetCol := (y - 1) % 3
	hasDot := false
	for i := targetRow * 3; i < targetRow*3+3; i++ {
		for j := targetCol * 3; j < targetCol*3+3; j++ {
			if board[i][j] == '.' {
				hasDot = true
			}
		}
	}
	if hasDot {
		for i := targetRow * 3; i < targetRow*3+3; i++ {
			for j := targetCol * 3; j < targetCol*3+3; j++ {
				if board[i][j] == '.' {
					board[i][j] = '!'
				}
			}
		}
	} else {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == '.' {
					board[i][j] = '!'
				}
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			writer.WriteByte(board[i][j])
			if j%3 == 2 && j != 8 {
				writer.WriteByte(' ')
			}
		}
		writer.WriteByte('\n')
		if i%3 == 2 && i != 8 {
			writer.WriteByte('\n')
		}
	}
	writer.Flush()
}
