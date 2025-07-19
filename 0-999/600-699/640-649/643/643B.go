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

	var n, k, a, b, c, d int
	if _, err := fmt.Fscan(reader, &n, &k, &a, &b, &c, &d); err != nil {
		return
	}
	if k <= n || n <= 4 {
		fmt.Fprintln(writer, -1)
		return
	}
	for i := 0; i < 2; i++ {
		fmt.Fprint(writer, a, " ", c, " ")
		for j := 1; j <= n; j++ {
			if j != a && j != b && j != c && j != d {
				fmt.Fprint(writer, j, " ")
			}
		}
		fmt.Fprintln(writer, d, b)
		a, c = c, a
		b, d = d, b
	}
}
