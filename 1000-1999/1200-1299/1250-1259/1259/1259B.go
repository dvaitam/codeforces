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
		seen := make(map[int]bool)
		ans := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			for x%2 == 0 {
				if !seen[x] {
					seen[x] = true
					ans++
				}
				x /= 2
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
