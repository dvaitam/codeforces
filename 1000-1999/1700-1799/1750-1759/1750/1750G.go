package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the actual algorithm described in problemG.txt.
// Currently this is just a placeholder that prints zero values.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	// Placeholder output
	for k := 1; k <= n; k++ {
		if k > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, 0)
	}
	fmt.Fprintln(writer)
}
