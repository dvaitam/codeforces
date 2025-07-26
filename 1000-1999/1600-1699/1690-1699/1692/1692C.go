package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		board := make([]string, 8)
		for i := 0; i < 8; i++ {
			fmt.Fscan(reader, &board[i])
		}
		found := false
		for i := 1; i < 7 && !found; i++ {
			for j := 1; j < 7 && !found; j++ {
				if board[i][j] == '#' &&
					board[i-1][j-1] == '#' &&
					board[i-1][j+1] == '#' &&
					board[i+1][j-1] == '#' &&
					board[i+1][j+1] == '#' {
					fmt.Fprintln(writer, i+1, j+1)
					found = true
				}
			}
		}
	}
}
