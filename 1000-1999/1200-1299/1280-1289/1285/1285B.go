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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		total := int64(0)
		for _, v := range a {
			total += v
		}

		curr := int64(0)
		maxPref := int64(-1 << 63)
		for i := 0; i < n-1; i++ {
			curr += a[i]
			if curr > maxPref {
				maxPref = curr
			}
			if curr < 0 {
				curr = 0
			}
		}

		curr = 0
		maxSuff := int64(-1 << 63)
		for i := n - 1; i >= 1; i-- {
			curr += a[i]
			if curr > maxSuff {
				maxSuff = curr
			}
			if curr < 0 {
				curr = 0
			}
		}

		if total > maxPref && total > maxSuff {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
