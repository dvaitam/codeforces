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
		var n, k int
		var s string
		fmt.Fscan(in, &n, &k)
		fmt.Fscan(in, &s)

		prev := -1
		ans := 0

		for i, ch := range s {
			if ch == '1' {
				if prev == -1 || i-prev >= k {
					ans++
				}
				prev = i
			}
		}

		fmt.Fprintln(out, ans)
	}
}
