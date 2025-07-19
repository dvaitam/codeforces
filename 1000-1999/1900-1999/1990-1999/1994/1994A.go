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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for tc := 0; tc < T; tc++ {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		if n == 1 && m == 1 {
			var a int
			fmt.Fscan(reader, &a)
			fmt.Fprintln(writer, -1)
			continue
		}
		total := n * m
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var a int
				fmt.Fscan(reader, &a)
				x := (a % total) + 1
				if j == m-1 {
					fmt.Fprintln(writer, x)
				} else {
					fmt.Fprint(writer, x, " ")
				}
			}
		}
	}
}
