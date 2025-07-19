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
		var l, r int64
		fmt.Fscan(reader, &n, &l, &r)
		ans := make([]int64, n)
		ok := true
		for i := 1; i <= n; i++ {
			ii := int64(i)
			x := ((l-1)/ii + 1) * ii
			if x > r {
				ok = false
				break
			}
			ans[i-1] = x
		}
		if ok {
			fmt.Fprintln(writer, "YES")
			for i, v := range ans {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v)
			}
			writer.WriteByte('\n')
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
