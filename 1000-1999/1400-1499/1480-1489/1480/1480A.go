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
		b := []byte(s)
		for i := 0; i < len(b); i++ {
			if i%2 == 0 {
				if b[i] == 'a' {
					b[i] = 'b'
				} else {
					b[i] = 'a'
				}
			} else {
				if b[i] == 'z' {
					b[i] = 'y'
				} else {
					b[i] = 'z'
				}
			}
		}
		fmt.Fprintln(out, string(b))
	}
}
