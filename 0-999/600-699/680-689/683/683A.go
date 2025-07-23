package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var a, x, y int
	if _, err := fmt.Fscan(in, &a, &x, &y); err != nil {
		return
	}

	if x > 0 && x < a && y > 0 && y < a {
		fmt.Fprintln(out, 0)
	} else if x >= 0 && x <= a && y >= 0 && y <= a && (x == 0 || x == a || y == 0 || y == a) {
		fmt.Fprintln(out, 1)
	} else {
		fmt.Fprintln(out, 2)
	}
}
