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
		var f, a, b int64
		fmt.Fscan(reader, &n, &f, &a, &b)
		times := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &times[i])
		}

		prev := int64(0)
		ok := true
		for i := 0; i < n; i++ {
			delta := times[i] - prev
			costOn := delta * a
			cost := costOn
			if b < cost {
				cost = b
			}
			f -= cost
			if f <= 0 {
				ok = false
				break
			}
			prev = times[i]
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
