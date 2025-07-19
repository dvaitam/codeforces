package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, x, y int64
	fmt.Fscan(reader, &n, &x, &y)
	if x == n && y == n {
		fmt.Fprint(os.Stdout, "Black")
	} else if max(abs(x-1), abs(y-1)) > max(abs(n-x), abs(n-y)) {
		fmt.Fprint(os.Stdout, "Black")
	} else {
		fmt.Fprint(os.Stdout, "White")
	}
}
