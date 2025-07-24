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
		total := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] == 2 {
				total++
			}
		}
		prefix := 0
		ans := -1
		if total%2 == 0 {
			for i := 0; i < n-1; i++ {
				if a[i] == 2 {
					prefix++
				}
				if prefix*2 == total {
					ans = i + 1
					break
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
