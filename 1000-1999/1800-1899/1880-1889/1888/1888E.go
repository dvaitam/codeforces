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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		hasValue := false
		var best int
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if !hasValue || v > best {
				best = v
				hasValue = true
			}
		}
		if hasValue {
			fmt.Fprintln(out, best)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}
