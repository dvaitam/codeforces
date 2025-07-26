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
		even, odd := 0, 0
		minVal := int(1<<31 - 1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x%2 == 0 {
				even++
			} else {
				odd++
			}
			if x < minVal {
				minVal = x
			}
		}
		if even == n || odd == n {
			fmt.Fprintln(writer, "YES")
		} else if minVal%2 == 1 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
