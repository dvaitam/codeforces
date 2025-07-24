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

	var n, m int
	fmt.Fscan(reader, &n, &m)
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	best := n + 1
	var bestPos []int
	for i := 0; i <= m-n; i++ {
		diff := make([]int, 0)
		for j := 0; j < n; j++ {
			if s[j] != t[i+j] {
				diff = append(diff, j+1)
			}
		}
		if len(diff) < best {
			best = len(diff)
			bestPos = diff
		}
	}

	fmt.Fprintln(writer, best)
	for i, p := range bestPos {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, p)
	}
	if len(bestPos) > 0 {
		fmt.Fprintln(writer)
	}
}
