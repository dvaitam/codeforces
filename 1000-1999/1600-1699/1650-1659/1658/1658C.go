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
		cnt1 := 0
		pos := -1
		for i, v := range c {
			if v == 1 {
				cnt1++
				pos = i
			}
		}
		if cnt1 != 1 {
			fmt.Fprintln(out, "NO")
			continue
		}
		ok := true
		for i := 1; i < n; i++ {
			prev := c[(pos+i-1)%n]
			cur := c[(pos+i)%n]
			if cur > prev+1 {
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
