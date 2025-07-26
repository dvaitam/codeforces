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
		var x int64
		var n int64
		fmt.Fscan(reader, &x, &n)
		r := n % 4
		var ans int64
		if x%2 == 0 {
			switch r {
			case 0:
				ans = x
			case 1:
				ans = x - n
			case 2:
				ans = x + 1
			case 3:
				ans = x + n + 1
			}
		} else {
			switch r {
			case 0:
				ans = x
			case 1:
				ans = x + n
			case 2:
				ans = x - 1
			case 3:
				ans = x - n - 1
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
