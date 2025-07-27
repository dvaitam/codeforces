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
		var n, m int64
		fmt.Fscan(reader, &n, &m)
		var ans int64
		if n == 1 {
			ans = 0
		} else if n == 2 {
			ans = m
		} else {
			ans = m * 2
		}
		fmt.Fprintln(writer, ans)
	}
}
