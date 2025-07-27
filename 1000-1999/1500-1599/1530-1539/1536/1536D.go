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
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		// remove consecutive duplicates
		c := []int{}
		for _, v := range b {
			if len(c) == 0 || c[len(c)-1] != v {
				c = append(c, v)
			}
		}
		inc := true
		dec := true
		for i := 1; i < len(c); i++ {
			if c[i] < c[i-1] {
				inc = false
			}
			if c[i] > c[i-1] {
				dec = false
			}
		}
		if inc || dec {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
