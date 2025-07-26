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
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		queue := make([]int, 0)
		first := c[0] + 1
		if first <= n {
			queue = append(queue, first)
		}
		ptr := 0
		a := make([]int, n)
		for i := 1; i <= n; i++ {
			for ptr < len(queue) && queue[ptr] < i {
				ptr++
			}
			if ptr < len(queue) && queue[ptr] == i {
				a[i-1] = 0
				ptr++
			} else {
				a[i-1] = 1
			}
			f := c[i-1] - (i-1)*a[i-1]
			if f < 0 {
				a[i-1] = 0
				f = c[i-1]
			}
			z := f + i
			if z > i && z <= n {
				queue = append(queue, z)
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		fmt.Fprintln(out)
	}
}
