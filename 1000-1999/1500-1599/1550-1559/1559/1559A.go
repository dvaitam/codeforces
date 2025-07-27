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
		var ans int
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if i == 0 {
				ans = x
			} else {
				ans &= x
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
