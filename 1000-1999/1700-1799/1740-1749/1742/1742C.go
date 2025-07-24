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
		ans := 'B'
		for i := 0; i < 8; i++ {
			row := board[i]
			allR := true
			for j := 0; j < 8; j++ {
				if row[j] != 'R' {
					allR = false
					break
				}
			}
			if allR {
				ans = 'R'
				break
			}
		}
		fmt.Fprintln(writer, string(ans))
	}
}
