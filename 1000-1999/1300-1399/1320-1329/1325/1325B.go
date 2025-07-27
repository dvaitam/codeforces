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
		var n int
		fmt.Fscan(reader, &n)
		m := make(map[int]struct{})
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			m[x] = struct{}{}
		}
		fmt.Fprintln(writer, len(m))
	}
}
