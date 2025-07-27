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
		ans := 0
		for i := 0; i < n; i++ {
			var r int
			fmt.Fscan(reader, &r)
			if r == 1 || r == 3 {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
