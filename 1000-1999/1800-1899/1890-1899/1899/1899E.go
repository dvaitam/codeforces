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

		minVal := a[0]
		pos := 0
		for i := 1; i < n; i++ {
			if a[i] < minVal {
				minVal = a[i]
				pos = i
			}
		}

		ok := true
		for i := pos + 1; i < n; i++ {
			if a[i] < a[i-1] {
				ok = false
				break
			}
		}

		if ok {
			fmt.Fprintln(out, pos)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
