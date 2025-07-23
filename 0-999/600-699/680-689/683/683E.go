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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	order := make([]int, 0, n)
	for id := n; id >= 1; id-- {
		pos := a[id]
		if pos < 0 {
			pos = 0
		} else if pos > len(order) {
			pos = len(order)
		}
		order = append(order, 0)
		copy(order[pos+1:], order[pos:])
		order[pos] = id
	}
	for i, v := range order {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
}
