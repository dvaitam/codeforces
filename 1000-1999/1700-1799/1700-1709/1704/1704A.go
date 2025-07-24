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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var a, b string
		fmt.Fscan(reader, &a)
		fmt.Fscan(reader, &b)

		ok := true
		// Check that the suffix of length m-1 of a matches that of b
		for i := 1; i < m && ok; i++ {
			if a[n-m+i] != b[i] {
				ok = false
			}
		}

		if ok {
			target := b[0]
			prefix := a[:n-m+1]
			found := false
			for i := 0; i < len(prefix); i++ {
				if prefix[i] == target {
					found = true
					break
				}
			}
			if !found {
				ok = false
			}
		}

		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
