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
		var n, x, y int64
		fmt.Fscan(in, &n, &x, &y)

		for i := int64(0); i < x; i++ {
			var tmp int64
			fmt.Fscan(in, &tmp)
		}

		k := x + y
		if k > n {
			k = n
		}
		ans := k - 2
		if ans < 0 {
			ans = 0
		}

		fmt.Fprintln(out, ans)
	}
}
