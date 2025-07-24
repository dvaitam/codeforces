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
	enter := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &enter[i])
	}
	exit := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &exit[i])
	}

	exitIndex := make(map[int]int, n)
	for i, v := range exit {
		exitIndex[v] = i
	}

	maxPos := -1
	fined := 0
	for _, car := range enter {
		pos := exitIndex[car]
		if pos > maxPos {
			maxPos = pos
		} else {
			fined++
		}
	}

	fmt.Fprintln(writer, fined)
}
