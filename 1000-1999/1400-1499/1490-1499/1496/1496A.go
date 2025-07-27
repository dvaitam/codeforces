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
		var n, k int
		var s string
		fmt.Fscan(reader, &n, &k)
		fmt.Fscan(reader, &s)
		if n < 2*k+1 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		ok := true
		for i := 0; i < k; i++ {
			if s[i] != s[n-1-i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
