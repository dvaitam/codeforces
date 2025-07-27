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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		ans := ((n + 1) / 2) * 5
		if ans < 15 {
			ans = 15
		}
		fmt.Fprintln(out, ans)
	}
}
