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
		matrix := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			matrix[i] = []byte(s)
		}

		operations := 0
		for i := 0; i < n/2; i++ {
			for j := 0; j < n/2; j++ {
				a := matrix[i][j]
				b := matrix[j][n-1-i]
				c := matrix[n-1-i][n-1-j]
				d := matrix[n-1-j][i]
				maxChar := a
				if b > maxChar {
					maxChar = b
				}
				if c > maxChar {
					maxChar = c
				}
				if d > maxChar {
					maxChar = d
				}
				operations += int(maxChar-a) + int(maxChar-b) + int(maxChar-c) + int(maxChar-d)
			}
		}
		fmt.Fprintln(out, operations)
	}
}
