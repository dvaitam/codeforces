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
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		pos[x] = i
	}

	best := 1
	cur := 1
	for i := 2; i <= n; i++ {
		if pos[i] > pos[i-1] {
			cur++
		} else {
			if cur > best {
				best = cur
			}
			cur = 1
		}
	}
	if cur > best {
		best = cur
	}

	fmt.Fprintln(writer, n-best)
}
