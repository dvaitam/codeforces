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

	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}

	present := make([]bool, 201)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if v >= 0 && v < len(present) {
			present[v] = true
		}
	}

	ops := 0
	for i := 0; i < x; i++ {
		if !present[i] {
			ops++
		}
	}
	if x >= 0 && x < len(present) && present[x] {
		ops++
	}

	fmt.Fprintln(writer, ops)
}
