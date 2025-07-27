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
		var s string
		fmt.Fscan(in, &s)
		maxDigit := byte('0')
		for i := 0; i < len(s); i++ {
			if s[i] > maxDigit {
				maxDigit = s[i]
			}
		}
		fmt.Fprintln(out, int(maxDigit-'0'))
	}
}
