package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		a, b := 0, 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				b++
			} else if s[i] == '0' && (i == 0 || s[i-1] != '0') {
				a++
			}
		}
		if b > a {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
