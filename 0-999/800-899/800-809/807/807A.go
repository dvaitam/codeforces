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

	ratings := make([]int, n)
	rated := false
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if a != b {
			rated = true
		}
		ratings[i] = a
	}

	if rated {
		fmt.Fprintln(writer, "rated")
		return
	}

	for i := 1; i < n; i++ {
		if ratings[i] > ratings[i-1] {
			fmt.Fprintln(writer, "unrated")
			return
		}
	}

	fmt.Fprintln(writer, "maybe")
}
