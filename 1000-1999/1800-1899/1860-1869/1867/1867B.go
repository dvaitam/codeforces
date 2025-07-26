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
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		mism := 0
		for i := 0; i < n/2; i++ {
			if s[i] != s[n-1-i] {
				mism++
			}
		}
		res := make([]byte, n+1)
		if n%2 == 1 {
			for i := mism; i <= n-mism; i++ {
				res[i] = '1'
			}
		} else {
			for i := mism; i <= n-mism; i += 2 {
				res[i] = '1'
			}
		}
		for i := 0; i <= n; i++ {
			if res[i] == 0 {
				res[i] = '0'
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
