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
		deg := make([]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			deg[u]++
			deg[v]++
		}
		leaves := 0
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				leaves++
			}
		}
		ans := (leaves + 1) / 2
		fmt.Fprintln(writer, ans)
	}
}
