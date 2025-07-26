package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		pos := make([]int, n)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			pos[v] = i
		}
		L, R := pos[0], pos[0]
		ans := int64(1)
		for x := 1; x < n; x++ {
			p := pos[x]
			if p < L {
				L = p
			} else if p > R {
				R = p
			} else {
				choices := int64(R - L + 1 - x)
				ans = (ans * choices) % mod
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
