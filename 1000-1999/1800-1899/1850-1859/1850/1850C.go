package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		grid := make([]string, 8)
		for i := 0; i < 8; i++ {
			fmt.Fscan(in, &grid[i])
		}
		var ans []byte
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				ch := grid[i][j]
				if ch != '.' {
					ans = append(ans, ch)
				}
			}
		}
		fmt.Fprintln(out, string(ans))
	}
}
