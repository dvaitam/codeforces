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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		k %= 4
		if k == 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		ok := true
		if k != 2 {
			for i := 1; i <= n; i += 2 {
				fmt.Fprintln(writer, i, i+1)
			}
		} else {
			for i := 1; i <= n; i += 2 {
				if ok {
					fmt.Fprintln(writer, i+1, i)
				} else {
					fmt.Fprintln(writer, i, i+1)
				}
				ok = !ok
			}
		}
	}
}
