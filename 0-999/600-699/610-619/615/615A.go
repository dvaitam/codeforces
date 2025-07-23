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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	bulbs := make([]bool, m+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		for j := 0; j < x; j++ {
			var y int
			fmt.Fscan(reader, &y)
			if y >= 1 && y <= m {
				bulbs[y] = true
			}
		}
	}
	for b := 1; b <= m; b++ {
		if !bulbs[b] {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	fmt.Fprintln(writer, "YES")
}
