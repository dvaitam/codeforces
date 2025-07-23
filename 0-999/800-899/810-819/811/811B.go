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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	for i := 0; i < m; i++ {
		var l, r, x int
		fmt.Fscan(in, &l, &r, &x)
		l--
		r--
		x--
		target := p[x]
		cnt := 0
		for j := l; j <= r; j++ {
			if p[j] < target {
				cnt++
			}
		}
		if l+cnt == x {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
