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
		var k int64
		fmt.Fscan(in, &n, &k)
		l := make([]int64, n)
		r := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &r[i])
		}
		var pref int64
		ans := int64(-1)
		for i := 0; i < n; i++ {
			seg := r[i] - l[i] + 1
			if pref+seg >= k {
				q := l[i] + (k - pref) - 1
				ans = q + 2
				break
			}
			pref += seg
		}
		fmt.Fprintln(out, ans)
	}
}
