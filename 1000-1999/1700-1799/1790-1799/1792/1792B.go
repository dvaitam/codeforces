package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a1, a2, a3, a4 int
		fmt.Fscan(reader, &a1, &a2, &a3, &a4)
		if a1 == 0 {
			fmt.Fprintln(writer, 1)
			continue
		}
		pair := min(a2, a3)
		ans := a1 + 2*pair
		diff := abs(a2 - a3)
		extra := diff + a4
		if a1+1 < extra {
			ans += a1 + 1
		} else {
			ans += extra
		}
		fmt.Fprintln(writer, ans)
	}
}
