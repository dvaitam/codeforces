package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		a := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Ints(a)
		ok := true
		for i := 0; i < n; i++ {
			if a[i+n]-a[i] < x {
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
