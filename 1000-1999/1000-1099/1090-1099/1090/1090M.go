package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	best := 0
	start := 0
	prev := -1

	for i := 0; i < n; i++ {
		var color int
		fmt.Fscan(in, &color)

		if i > 0 && color == prev {
			start = i
		}
		length := i - start + 1
		if length > best {
			best = length
		}
		prev = color
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, best)
	out.Flush()
}
