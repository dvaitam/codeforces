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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		board := make([]string, 3)
		for i := 0; i < 3; i++ {
			fmt.Fscan(reader, &board[i])
		}
		winner := checkWinner(board)
		if winner == 0 {
			fmt.Fprintln(writer, "DRAW")
		} else {
			fmt.Fprintf(writer, "%c\n", winner)
		}
	}
}

func checkWinner(b []string) byte {
	lines := [][][2]int{{{0, 0}, {0, 1}, {0, 2}}, {{1, 0}, {1, 1}, {1, 2}}, {{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}}, {{0, 1}, {1, 1}, {2, 1}}, {{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}}, {{0, 2}, {1, 1}, {2, 0}}}
	for _, line := range lines {
		c1 := b[line[0][0]][line[0][1]]
		if c1 == '.' {
			continue
		}
		if b[line[1][0]][line[1][1]] == c1 && b[line[2][0]][line[2][1]] == c1 {
			return c1
		}
	}
	return 0
}
