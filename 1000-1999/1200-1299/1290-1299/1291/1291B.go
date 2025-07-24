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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		prefix := make([]bool, n)
		ok := true
		for i := 0; i < n; i++ {
			if ok && a[i] >= int64(i) {
				prefix[i] = true
			} else {
				ok = false
			}
		}
		suffix := make([]bool, n)
		ok = true
		for i := n - 1; i >= 0; i-- {
			if ok && a[i] >= int64(n-1-i) {
				suffix[i] = true
			} else {
				ok = false
			}
		}
		possible := false
		for i := 0; i < n; i++ {
			if prefix[i] && suffix[i] {
				possible = true
				break
			}
		}
		if possible {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
