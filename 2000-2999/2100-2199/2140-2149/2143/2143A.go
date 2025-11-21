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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if n <= 2 {
			fmt.Fprintln(out, "YES")
			continue
		}
		mode := 0 // 0 for increasing, 1 for decreasing
		ok := true
		for i := 1; i < n; i++ {
			if mode == 0 {
				if p[i] > p[i-1] {
					continue
				}
				if p[i] < p[i-1] {
					mode = 1
					continue
				}
				ok = false
				break
			} else {
				if p[i] < p[i-1] {
					continue
				}
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
