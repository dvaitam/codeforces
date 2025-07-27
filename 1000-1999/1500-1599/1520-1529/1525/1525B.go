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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		sorted := true
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] != i+1 {
				sorted = false
			}
		}
		if sorted {
			fmt.Fprintln(out, 0)
			continue
		}
		if a[0] == 1 || a[n-1] == n {
			fmt.Fprintln(out, 1)
		} else if a[0] == n && a[n-1] == 1 {
			fmt.Fprintln(out, 3)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
