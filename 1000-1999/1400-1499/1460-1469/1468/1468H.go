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
		var n, k, m int
		fmt.Fscan(reader, &n, &k, &m)
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if (n-m)%(k-1) != 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		need := (k - 1) / 2
		ok := false
		for i := 0; i < m; i++ {
			left := b[i] - (i + 1)
			right := n - b[i] - (m - i - 1)
			if left >= need && right >= need {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
