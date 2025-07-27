package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	maxNode := -1
	for {
		var u, v int
		if _, err := fmt.Fscan(r, &u, &v); err != nil {
			break
		}
		if u > maxNode {
			maxNode = u
		}
		if v > maxNode {
			maxNode = v
		}
	}
	n := maxNode + 1

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for i := 0; i < n; i++ {
		fmt.Fprintln(w, i)
	}
}
