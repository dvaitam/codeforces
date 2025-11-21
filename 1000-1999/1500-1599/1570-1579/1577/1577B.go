package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	maxVal := int64(-1 << 63)
	for i := 0; i < n; i++ {
		var x int64
		if _, err := fmt.Fscan(in, &x); err != nil {
			return
		}
		if x > maxVal {
			maxVal = x
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, maxVal)
}
