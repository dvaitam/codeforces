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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	last := make(map[int]int)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(in, &x)
		last[x] = i
	}

	res := 0
	minPos := n + 1
	for cafe, pos := range last {
		if pos < minPos {
			minPos = pos
			res = cafe
		}
	}

	fmt.Fprintln(out, res)
}
