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
		cnt1, cnt2 := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == 1 {
				cnt1++
			} else if x == 2 {
				cnt2++
			}
		}
		if cnt1%2 != 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if cnt1 == 0 && cnt2%2 == 1 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
	}
}
