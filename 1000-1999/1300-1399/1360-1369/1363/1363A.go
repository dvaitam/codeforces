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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(reader, &n, &x)
		odd := 0
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			if v%2 != 0 {
				odd++
			}
		}
		even := n - odd
		possible := false
		for k := 1; k <= odd && k <= x; k += 2 {
			if x-k <= even {
				possible = true
				break
			}
		}
		if possible {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
