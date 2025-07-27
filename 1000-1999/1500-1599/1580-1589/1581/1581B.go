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
		var n, m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		maxEdges := n * (n - 1) / 2
		if m < n-1 || m > maxEdges {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if n == 1 {
			if m == 0 && k > 1 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
			continue
		}
		if n == 2 {
			if m == 1 && k > 2 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
			continue
		}
		var minDiameter int64
		if m == maxEdges {
			minDiameter = 1
		} else {
			minDiameter = 2
		}
		if k-1 > minDiameter {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
