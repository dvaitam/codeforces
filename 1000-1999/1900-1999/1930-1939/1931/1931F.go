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
		var n, k int
		fmt.Fscan(in, &n, &k)

		base := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &base[i])
		}
		pos := make([]int, n+1)
		for i, v := range base {
			pos[v] = i
		}

		arr := make([]int, n)
		ok := true
		for s := 1; s < k; s++ {
			for i := 0; i < n; i++ {
				fmt.Fscan(in, &arr[i])
			}
			start := pos[arr[0]]
			for i := 0; i < n && ok; i++ {
				if arr[i] != base[(start+i)%n] {
					ok = false
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
