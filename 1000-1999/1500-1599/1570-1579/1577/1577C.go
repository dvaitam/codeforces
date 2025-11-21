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

	isSorted := true
	var prev int64
	for i := 0; i < n; i++ {
		var x int64
		if _, err := fmt.Fscan(in, &x); err != nil {
			return
		}
		if i > 0 && x < prev {
			isSorted = false
		}
		prev = x
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	if isSorted {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
