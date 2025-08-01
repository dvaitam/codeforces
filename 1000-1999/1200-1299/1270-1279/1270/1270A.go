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
		var n, k1, k2 int
		fmt.Fscan(reader, &n, &k1, &k2)
		hasMax := false
		for i := 0; i < k1; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == n {
				hasMax = true
			}
		}
		for i := 0; i < k2; i++ {
			var x int
			fmt.Fscan(reader, &x)
		}
		if hasMax {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
