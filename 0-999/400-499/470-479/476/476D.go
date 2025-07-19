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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	// minimal m is (6*n - 1) * k
	m := (6*n - 1) * k
	fmt.Fprintln(writer, m)
	for i := 0; i < n; i++ {
		x := 1 + 6*i
		a := x * k
		b := (x + 2) * k
		c := (x + 4) * k
		var d int
		if (x+1)%3 != 0 {
			d = (x + 1) * k
		} else {
			d = (x + 3) * k
		}
		fmt.Fprintln(writer, a, b, c, d)
	}
}
