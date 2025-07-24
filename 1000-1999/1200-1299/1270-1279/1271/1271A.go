package main

import (
	"bufio"
	"fmt"
	"os"
)

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

	var a, b, c, d, e, f int
	if _, err := fmt.Fscan(reader, &a, &b, &c, &d, &e, &f); err != nil {
		return
	}

	var res int
	if e > f {
		t1 := min(a, d)
		res += t1 * e
		d -= t1
		t2 := min(min(b, c), d)
		res += t2 * f
	} else {
		t2 := min(min(b, c), d)
		res += t2 * f
		d -= t2
		t1 := min(a, d)
		res += t1 * e
	}

	fmt.Fprintln(writer, res)
}
