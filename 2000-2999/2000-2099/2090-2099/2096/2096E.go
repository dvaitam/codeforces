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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		res := int64(0)
		countB := 0
		for _, ch := range s {
			if ch == 'B' {
				res += int64(countB)
			} else {
				countB++
			}
		}
		fmt.Fprintln(out, res)
	}
}
