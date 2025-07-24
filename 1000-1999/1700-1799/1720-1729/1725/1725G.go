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

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	var ans int64
	if n == 1 {
		ans = 3
	} else {
		k := n - 2
		q := k / 3
		r := k % 3
		ans = 5 + 4*q
		if r == 1 {
			ans += 2
		} else if r == 2 {
			ans += 3
		}
	}

	fmt.Fprintln(out, ans)
}
