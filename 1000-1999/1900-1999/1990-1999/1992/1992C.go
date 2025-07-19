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
	for t > 0 {
		t--
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		for i := n; i >= k; i-- {
			fmt.Fprint(writer, i, " ")
		}
		for i := m + 1; i < k; i++ {
			fmt.Fprint(writer, i, " ")
		}
		for i := 1; i <= m; i++ {
			fmt.Fprint(writer, i, " ")
		}
		fmt.Fprint(writer, "\n")
	}
}
