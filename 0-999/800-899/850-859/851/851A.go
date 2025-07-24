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

	var n, k, t int
	if _, err := fmt.Fscan(in, &n, &k, &t); err != nil {
		return
	}

	var ans int
	switch {
	case t <= k:
		ans = t
	case t <= n:
		ans = k
	default:
		ans = n + k - t
	}

	if ans < 0 {
		ans = 0
	}
	fmt.Fprintln(out, ans)
}
