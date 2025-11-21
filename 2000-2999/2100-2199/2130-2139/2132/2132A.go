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
		var a string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &a)
		var m int
		var b, c string
		fmt.Fscan(in, &m)
		fmt.Fscan(in, &b)
		fmt.Fscan(in, &c)

		left := make([]byte, 0, m)
		right := []byte(a)
		for i := 0; i < m; i++ {
			if c[i] == 'V' {
				left = append([]byte{b[i]}, left...)
			} else {
				right = append(right, b[i])
			}
		}

		fmt.Fprintln(out, string(append(left, right...)))
	}
}
