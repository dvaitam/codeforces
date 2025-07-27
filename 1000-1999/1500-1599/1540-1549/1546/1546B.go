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
		var n, m int
		fmt.Fscan(in, &n, &m)
		res := make([]byte, m)
		var s string
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s)
			for j := 0; j < m; j++ {
				res[j] ^= s[j]
			}
		}
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &s)
			for j := 0; j < m; j++ {
				res[j] ^= s[j]
			}
		}
		out.Write(res)
		out.WriteByte('\n')
	}
}
