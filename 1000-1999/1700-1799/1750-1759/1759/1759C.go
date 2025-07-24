package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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
		var l, r, x int
		fmt.Fscan(reader, &l, &r, &x)
		var a, b int
		fmt.Fscan(reader, &a, &b)

		ans := -1
		if a == b {
			ans = 0
		} else if abs(a-b) >= x {
			ans = 1
		} else if abs(b-l) < x && abs(b-r) < x {
			ans = -1
		} else if abs(a-l) >= x && abs(b-l) >= x {
			ans = 2
		} else if abs(a-r) >= x && abs(b-r) >= x {
			ans = 2
		} else if abs(r-l) >= x && ((abs(a-l) >= x && abs(r-b) >= x) || (abs(a-r) >= x && abs(l-b) >= x)) {
			ans = 3
		} else {
			ans = -1
		}
		fmt.Fprintln(writer, ans)
	}
}
