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
		ans := 1<<31 - 1
		for i := 0; i < n; i++ {
			var d, s int
			fmt.Fscan(reader, &d, &s)
			val := d + (s-1)/2
			if val < ans {
				ans = val
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
