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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(writer, 1)
			continue
		}
		alt := a[0] != a[1]
		if alt {
			for i := 2; i < n; i++ {
				if a[i] != a[i%2] {
					alt = false
					break
				}
			}
		}
		if alt {
			fmt.Fprintln(writer, (n+3)/2)
		} else {
			fmt.Fprintln(writer, n)
		}
	}
}
