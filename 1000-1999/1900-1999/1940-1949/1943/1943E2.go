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
		var m int
		var k int64
		fmt.Fscan(in, &m, &k)
		freq := make([]int64, m+1)
		for i := 0; i <= m; i++ {
			fmt.Fscan(in, &freq[i])
		}
		ans := 0
		for i := 0; ; i++ {
			var cnt int64
			if i <= m {
				cnt = freq[i]
			}
			if cnt <= int64(i)*k {
				ans = i
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
