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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var x int
		fmt.Fscan(reader, &n, &x)
		stacks := make([][]int, 3)
		for s := 0; s < 3; s++ {
			stacks[s] = make([]int, n)
			for i := 0; i < n; i++ {
				fmt.Fscan(reader, &stacks[s][i])
			}
		}
		cur := 0
		for s := 0; s < 3 && cur != x; s++ {
			for i := 0; i < n && cur != x; i++ {
				v := stacks[s][i]
				if v|x != x {
					break
				}
				cur |= v
			}
		}
		if cur == x {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
