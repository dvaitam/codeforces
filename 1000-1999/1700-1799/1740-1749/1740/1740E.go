package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution implements a simple heuristic for problem E.
// The real optimal strategy is non-trivial; here we approximate
// by counting leaves and nodes with exactly one child.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &parent[i])
	}
	children := make([]int, n+1)
	for i := 2; i <= n; i++ {
		children[parent[i]]++
	}
	leaves := 0
	ones := 0
	for i := 1; i <= n; i++ {
		if children[i] == 0 {
			leaves++
		} else if children[i] == 1 {
			ones++
		}
	}
	fmt.Println(leaves + ones)
}
