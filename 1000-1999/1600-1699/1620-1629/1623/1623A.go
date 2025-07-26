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
		var n, m, rb, cb, rd, cd int
		fmt.Fscan(reader, &n, &m, &rb, &cb, &rd, &cd)
		r, c := rb, cb
		dr, dc := 1, 1
		steps := 0
		for {
			if r == rd || c == cd {
				fmt.Fprintln(writer, steps)
				break
			}
			if r+dr > n || r+dr < 1 {
				dr = -dr
			}
			if c+dc > m || c+dc < 1 {
				dc = -dc
			}
			r += dr
			c += dc
			steps++
		}
	}
}
