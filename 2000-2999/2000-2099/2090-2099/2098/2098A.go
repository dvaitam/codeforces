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
		var s string
		fmt.Fscan(in, &s)
		cnt := make([]int, 10)
		for _, ch := range s {
			cnt[int(ch-'0')]++
		}
		res := make([]byte, 10)
		for i := 0; i < 10; i++ {
			thr := 9 - i
			for d := thr; d <= 9; d++ {
				if cnt[d] > 0 {
					res[i] = byte('0' + d)
					cnt[d]--
					break
				}
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
