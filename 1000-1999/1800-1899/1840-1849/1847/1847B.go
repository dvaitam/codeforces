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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		all := a[0]
		for i := 1; i < n; i++ {
			all &= a[i]
		}
		if all > 0 {
			fmt.Fprintln(writer, 1)
			continue
		}
		cur := -1
		cnt := 0
		for i := 0; i < n; i++ {
			cur &= a[i]
			if cur == 0 {
				cnt++
				cur = -1
			}
		}
		fmt.Fprintln(writer, cnt)
	}
}
