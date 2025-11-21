package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		result := 1
		for i := 0; i < m; i++ {
			result *= n
		}

		fmt.Fprintln(out, result)
	}
}

