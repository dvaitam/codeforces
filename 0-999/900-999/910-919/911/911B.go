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

	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}

	ans := 0
	for i := 1; i <= n-1; i++ {
		x := a / i
		y := b / (n - i)
		if x < y {
			y = x
		}
		if y > ans {
			ans = y
		}
	}

	fmt.Fprintln(out, ans)
}
