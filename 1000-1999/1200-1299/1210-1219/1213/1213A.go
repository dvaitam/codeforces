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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	even, odd := 0, 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x%2 == 0 {
			even++
		} else {
			odd++
		}
	}
	if even < odd {
		fmt.Fprintln(out, even)
	} else {
		fmt.Fprintln(out, odd)
	}
}
