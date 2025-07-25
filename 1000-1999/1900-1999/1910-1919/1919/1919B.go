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
		var s string
		fmt.Fscan(in, &s)
		diff := 0
		for i := 0; i < n; i++ {
			if s[i] == '+' {
				diff++
			} else {
				diff--
			}
		}
		if diff < 0 {
			diff = -diff
		}
		fmt.Fprintln(out, diff)
	}
}
