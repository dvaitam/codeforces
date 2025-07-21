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
		var a, b string
		fmt.Fscan(in, &a, &b)
		as := []byte(a)
		bs := []byte(b)
		n := len(as)
		m := len(bs)
		ans := n + m
		for i := 0; i < m; i++ {
			j := i
			for _, ch := range as {
				if j < m && ch == bs[j] {
					j++
				}
			}
			val := n + m - (j - i)
			if val < ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}
