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
		var n, r int64
		fmt.Fscan(reader, &n, &r)
		if r < n {
			ans := r * (r + 1) / 2
			fmt.Fprintln(writer, ans)
		} else {
			ans := n*(n-1)/2 + 1
			fmt.Fprintln(writer, ans)
		}
	}
}
