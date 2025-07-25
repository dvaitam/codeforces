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
		a := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		ok := true
		// check row sums
		for i := 0; i < n && ok; i++ {
			sum := 0
			for j := 0; j < m; j++ {
				ai := int(a[i][j] - '0')
				bi := int(b[i][j] - '0')
				diff := (bi - ai) % 3
				if diff < 0 {
					diff += 3
				}
				sum += diff
			}
			if sum%3 != 0 {
				ok = false
			}
		}
		// check column sums
		for j := 0; j < m && ok; j++ {
			sum := 0
			for i := 0; i < n; i++ {
				ai := int(a[i][j] - '0')
				bi := int(b[i][j] - '0')
				diff := (bi - ai) % 3
				if diff < 0 {
					diff += 3
				}
				sum += diff
			}
			if sum%3 != 0 {
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
