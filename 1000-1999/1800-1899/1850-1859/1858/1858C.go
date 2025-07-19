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

	var t, n int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		fmt.Fscan(reader, &n)
		// For each odd starting point, double until exceeding n
		for i := 1; i <= n; i += 2 {
			for j := i; j <= n; j <<= 1 {
				fmt.Fprint(writer, j, " ")
			}
		}
		fmt.Fprintln(writer)
	}
}
