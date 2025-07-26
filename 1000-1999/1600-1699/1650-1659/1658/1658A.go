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
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		ans := 0
		for i := 1; i < n; i++ {
			if s[i-1] == '0' && s[i] == '0' {
				ans += 2
			}
		}
		for i := 2; i < n; i++ {
			if s[i-2] == '0' && s[i-1] == '1' && s[i] == '0' {
				ans += 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
