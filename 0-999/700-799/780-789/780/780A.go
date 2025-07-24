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

	table := make(map[int]bool)
	current := 0
	maxCount := 0
	for i := 0; i < 2*n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if table[x] {
			delete(table, x)
			current--
		} else {
			table[x] = true
			current++
			if current > maxCount {
				maxCount = current
			}
		}
	}
	fmt.Fprintln(writer, maxCount)
}
