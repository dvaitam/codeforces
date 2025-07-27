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
		minVal := 101
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] < minVal {
				minVal = a[i]
			}
		}
		cnt := 0
		for i := 0; i < n; i++ {
			if a[i] > minVal {
				cnt++
			}
		}
		fmt.Fprintln(out, cnt)
	}
}
