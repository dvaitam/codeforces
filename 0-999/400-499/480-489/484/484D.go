package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	var ans, u, v int64
	const negInf int64 = -1e18
	u, v = negInf, negInf
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		if ans+x > u {
			u = ans + x
		}
		if ans-x > v {
			v = ans - x
		}
		ans = max(u-x, v+x)
	}
	fmt.Fprint(out, ans)
}
