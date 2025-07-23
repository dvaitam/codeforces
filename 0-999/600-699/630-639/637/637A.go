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
	counts := make(map[int]int)
	bestID := 0
	bestCount := 0
	for i := 0; i < n; i++ {
		var id int
		fmt.Fscan(reader, &id)
		counts[id]++
		if counts[id] > bestCount {
			bestCount = counts[id]
			bestID = id
		}
	}
	fmt.Fprintln(writer, bestID)
}
