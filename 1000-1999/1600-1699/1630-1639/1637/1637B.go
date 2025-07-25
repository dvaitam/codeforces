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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		totalLen := int64(n * (n + 1) * (n + 2) / 6)
		addZero := int64(0)
		for i := 0; i < n; i++ {
			if a[i] == 0 {
				left := i + 1
				right := n - i
				addZero += int64(left * right)
			}
		}
		fmt.Fprintln(out, totalLen+addZero)
	}
}
