package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]string, 2*n-2)
		for i := range arr {
			fmt.Fscan(reader, &arr[i])
		}
		var a, b string
		for _, s := range arr {
			if len(s) == n-1 {
				if a == "" {
					a = s
				} else {
					b = s
				}
			}
		}
		if reverse(a) == b {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
