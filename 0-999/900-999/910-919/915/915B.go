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

	var n, pos, l, r int
	if _, err := fmt.Fscan(reader, &n, &pos, &l, &r); err != nil {
		return
	}

	var ans int
	switch {
	case l == 1 && r == n:
		ans = 0
	case l == 1:
		ans = abs(pos-r) + 1
	case r == n:
		ans = abs(pos-l) + 1
	default:
		ans = min(abs(pos-l), abs(pos-r)) + (r - l) + 2
	}

	fmt.Fprintln(writer, ans)
}
