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
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		oddParity := arr[0] % 2
		evenParity := 0
		if n > 1 {
			evenParity = arr[1] % 2
		}
		ok := true
		for i := 0; i < n && ok; i++ {
			if i%2 == 0 {
				if arr[i]%2 != oddParity {
					ok = false
				}
			} else {
				if arr[i]%2 != evenParity {
					ok = false
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
