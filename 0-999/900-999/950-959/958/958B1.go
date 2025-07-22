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
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		deg[u]++
		deg[v]++
	}
	count := 0
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			count++
		}
	}
	fmt.Fprintln(writer, count)
}
