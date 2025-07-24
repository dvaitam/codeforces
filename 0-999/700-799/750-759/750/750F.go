package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var h int
		if _, err := fmt.Fscan(in, &h); err != nil {
			return
		}
		n := (1 << h) - 1
		deg := make([]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			deg[u]++
			deg[v]++
		}
		root := 1
		for i := 1; i <= n; i++ {
			if deg[i] == 2 {
				root = i
				break
			}
		}
		fmt.Println(root)
	}
}
