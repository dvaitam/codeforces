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
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		if a > b {
			fmt.Fprintln(writer, 1)
		} else {
			ans := (n + a - 1) / a
			if ans < 1 {
				ans = 1
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
