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
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s string
		fmt.Fscan(in, &s)
		cnt := make([]int, 7)
		for _, ch := range s {
			if ch >= 'A' && ch <= 'G' {
				cnt[ch-'A']++
			}
		}
		add := 0
		for i := 0; i < 7; i++ {
			if m > cnt[i] {
				add += m - cnt[i]
			}
		}
		fmt.Fprintln(out, add)
	}
}
