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
		idx := -1
		for i := 0; i+1 < n; i++ {
			if a[i] > a[i+1] {
				idx = i
				break
			}
		}
		if idx == -1 {
			fmt.Fprintln(out, "NO")
			continue
		}
		fmt.Fprintln(out, "YES")
		fmt.Fprintln(out, 2)
		fmt.Fprintf(out, "%d %d\n", a[idx], a[idx+1])
	}
}
