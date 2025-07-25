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
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}
		avg := sum / int64(n)
		var prefix int64
		ok := true
		for i := 0; i < n; i++ {
			prefix += a[i]
			if prefix < int64(i+1)*avg {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
