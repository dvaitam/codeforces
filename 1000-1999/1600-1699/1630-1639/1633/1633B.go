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
		var s string
		fmt.Fscan(in, &s)
		zeros := 0
		for _, ch := range s {
			if ch == '0' {
				zeros++
			}
		}
		ones := len(s) - zeros
		if zeros == ones {
			fmt.Fprintln(out, zeros-1)
		} else if zeros < ones {
			fmt.Fprintln(out, zeros)
		} else {
			fmt.Fprintln(out, ones)
		}
	}
}
