package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &perm[i])
		}
		target1 := n / 2
		target2 := target1 + 1
		pos1, pos2 := -1, -1
		for i, val := range perm {
			if val == target1 {
				pos1 = i + 1
			} else if val == target2 {
				pos2 = i + 1
			}
		}
		fmt.Fprintf(out, "%d %d\n", pos1, pos2)
	}
}
