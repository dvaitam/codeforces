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
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}

	ans := make([]int64, n)
	maxRight := int64(0)
	for i := n - 1; i >= 0; i-- {
		if h[i] > maxRight {
			ans[i] = 0
			maxRight = h[i]
		} else {
			ans[i] = maxRight + 1 - h[i]
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, ans[i])
	}
	out.WriteByte('\n')
}
