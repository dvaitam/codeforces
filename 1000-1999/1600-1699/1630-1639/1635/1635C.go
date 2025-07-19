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
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// If second last > last, impossible
		if a[n-1] > a[n] {
			fmt.Fprintln(writer, -1)
			continue
		}
		// If last element negative
		if a[n] < 0 {
			ok := true
			for i := 2; i <= n; i++ {
				if a[i] < a[i-1] {
					ok = false
					break
				}
			}
			if !ok {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, 0)
			}
			continue
		}
		// Otherwise, apply operations
		// We will make pairs (i, n-1, n) for i from n-2 down to 1
		fmt.Fprintln(writer, n-2)
		for i := n - 2; i >= 1; i-- {
			fmt.Fprintln(writer, i, n-1, n)
		}
	}
}
