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
		var n, a, b int64
		fmt.Fscan(reader, &n, &a, &b)
		if a == 1 {
			if (n-1)%b == 0 {
				fmt.Fprintln(writer, "Yes")
			} else {
				fmt.Fprintln(writer, "No")
			}
			continue
		}
		x := int64(1)
		ok := false
		for x <= n {
			if (n-x)%b == 0 {
				ok = true
				break
			}
			if x > n/a {
				break
			}
			x *= a
		}
		if ok {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
