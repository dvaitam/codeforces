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
		var n int
		fmt.Fscan(reader, &n)
		sum := int64(0)
		ones := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += int64(x)
			if x == 1 {
				ones++
			}
		}
		if n == 1 || int64(ones) > sum-int64(n) {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
