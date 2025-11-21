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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &rows[i])
		}
		for i := n - 1; i >= 0; i-- {
			for j := 0; j < 4; j++ {
				if rows[i][j] == '#' {
					fmt.Fprint(out, j+1)
					if i != 0 {
						fmt.Fprint(out, " ")
					}
					break
				}
			}
		}
		if t > 1 {
			fmt.Fprint(out, "\n")
		} else {
			fmt.Fprint(out, "\n")
		}
	}
}
