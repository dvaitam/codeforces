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
	counts := make([]int, 1001)
	maxHeight, distinct := 0, 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if counts[x] == 0 {
			distinct++
		}
		counts[x]++
		if counts[x] > maxHeight {
			maxHeight = counts[x]
		}
	}
	fmt.Fprintln(writer, maxHeight, distinct)
}
