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
		mx := int64(0)
		x := int64(0)
		ans := int64(0)
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(reader, &v)
			if diff := mx - v; diff > 0 {
				ans += diff
				if diff > x {
					x = diff
				}
			}
			if v > mx {
				mx = v
			}
		}
		fmt.Fprintln(writer, ans+x)
	}
}
