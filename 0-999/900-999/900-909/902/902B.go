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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &parent[i])
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &color[i])
	}

	steps := 1 // color root
	for i := 2; i <= n; i++ {
		if color[i] != color[parent[i]] {
			steps++
		}
	}
	fmt.Fprintln(writer, steps)
}
