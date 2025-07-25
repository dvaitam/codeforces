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

	var m int
	fmt.Fscan(reader, &m)

	counts := make([]int64, 30)
	for i := 0; i < m; i++ {
		var t, v int
		fmt.Fscan(reader, &t, &v)
		if t == 1 {
			if v >= 0 && v < 30 {
				counts[v]++
			}
		} else {
			w := int64(v)
			carry := int64(0)
			possible := true
			for b := 0; b < 30; b++ {
				carry += counts[b]
				if (w>>b)&1 == 1 {
					if carry == 0 {
						possible = false
						break
					}
					carry--
				}
				carry /= 2
			}
			if possible {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
