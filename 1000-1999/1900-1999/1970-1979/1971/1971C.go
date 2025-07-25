package main

import (
	"bufio"
	"fmt"
	"os"
)

func between(a, b, x int) bool {
	if x == a || x == b {
		return false
	}
	cur := a
	for {
		cur = cur%12 + 1
		if cur == x {
			return true
		}
		if cur == b {
			break
		}
	}
	return false
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
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b, &c, &d)
		abC := between(a, b, c)
		abD := between(a, b, d)
		cdA := between(c, d, a)
		cdB := between(c, d, b)
		if abC != abD && cdA != cdB {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
