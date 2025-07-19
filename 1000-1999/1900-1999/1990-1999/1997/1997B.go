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
		fmt.Fscan(reader, &n)
		var a, b string
		fmt.Fscan(reader, &a, &b)
		ans := 0
		for i := 0; i < n-2; i++ {
			if a[i] != b[i] && a[i] == a[i+2] && b[i] == b[i+2] && a[i+1] == '.' && b[i+1] == '.' {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
