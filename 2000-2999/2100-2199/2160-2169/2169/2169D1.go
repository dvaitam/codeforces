package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit int64 = 1_000_000_000_000

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var x int
		var y, k int64
		fmt.Fscan(in, &x, &y, &k)

		if y == 1 {
			fmt.Fprintln(out, -1)
			continue
		}

		cur := k
		ok := true
		for i := 0; i < x && cur <= limit; i++ {
			add := (cur - 1) / (y - 1)
			cur += add
			if cur > limit {
				ok = false
				break
			}
		}

		if ok {
			fmt.Fprintln(out, cur)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
