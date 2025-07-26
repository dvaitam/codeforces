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
		var altSum int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			if i%2 == 0 {
				altSum += x
			} else {
				altSum -= x
			}
		}
		if n%2 == 1 || altSum <= 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
