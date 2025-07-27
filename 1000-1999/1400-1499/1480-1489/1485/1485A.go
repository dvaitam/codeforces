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
		var a, b int
		fmt.Fscan(reader, &a, &b)
		ans := int(1e9)
		for inc := 0; inc <= 30; inc++ {
			nb := b + inc
			if nb == 1 {
				continue
			}
			cnt := inc
			x := a
			for x > 0 {
				x /= nb
				cnt++
			}
			if cnt < ans {
				ans = cnt
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
