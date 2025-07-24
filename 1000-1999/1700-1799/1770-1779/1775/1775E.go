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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		prefix := int64(0)
		maxPref := int64(0)
		minPref := int64(0)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			prefix += x
			if prefix > maxPref {
				maxPref = prefix
			}
			if prefix < minPref {
				minPref = prefix
			}
		}
		fmt.Fprintln(writer, maxPref-minPref)
	}
}
