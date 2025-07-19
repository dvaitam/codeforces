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
		allEqual := true
		for i := 1; i < n; i++ {
			if a[i] != a[0] {
				allEqual = false
				break
			}
		}
		if allEqual {
			writer.WriteString("-1\n")
			continue
		}
		mx := a[0]
		for _, v := range a {
			if v > mx {
				mx = v
			}
		}
		var b, c []int
		for _, v := range a {
			if v != mx {
				b = append(b, v)
			} else {
				c = append(c, v)
			}
		}
		// output sizes
		fmt.Fprintf(writer, "%d %d\n", len(b), len(c))
		// output b
		for _, v := range b {
			fmt.Fprintf(writer, "%d ", v)
		}
		writer.WriteByte('\n')
		// output c
		for _, v := range c {
			fmt.Fprintf(writer, "%d ", v)
		}
		writer.WriteByte('\n')
	}
}
