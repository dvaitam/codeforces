package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		cur := l
		diff := int64(1)
		length := int64(1)
		for {
			if cur+diff > r {
				break
			}
			cur += diff
			diff++
			length++
		}
		fmt.Fprintln(out, length)
	}
}
