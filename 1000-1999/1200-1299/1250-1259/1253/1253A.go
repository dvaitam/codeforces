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
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		diff := make([]int, n)
		possible := true
		for i := 0; i < n; i++ {
			diff[i] = b[i] - a[i]
			if diff[i] < 0 {
				possible = false
			}
		}
		if !possible {
			fmt.Fprintln(out, "NO")
			continue
		}
		i := 0
		for i < n && diff[i] == 0 {
			i++
		}
		if i == n {
			fmt.Fprintln(out, "YES")
			continue
		}
		k := diff[i]
		if k <= 0 {
			possible = false
		} else {
			for i < n && diff[i] == k {
				i++
			}
			for i < n {
				if diff[i] != 0 {
					possible = false
					break
				}
				i++
			}
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
