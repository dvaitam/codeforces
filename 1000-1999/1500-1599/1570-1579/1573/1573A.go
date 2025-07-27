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
		var s string
		fmt.Fscan(in, &s)
		ans := 0
		for i := 0; i < n; i++ {
			d := int(s[i] - '0')
			ans += d
			if i < n-1 && d > 0 {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
