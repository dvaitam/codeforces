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
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	words := make(map[string]bool, n)
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(reader, &w)
		words[w] = true
	}

	common := 0
	for i := 0; i < m; i++ {
		var w string
		fmt.Fscan(reader, &w)
		if words[w] {
			common++
		}
	}

	if common%2 == 1 {
		n++
	}
	if n > m {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
