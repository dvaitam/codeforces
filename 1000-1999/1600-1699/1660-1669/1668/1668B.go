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
		var m int64
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		var sum, mn, mx int64
		mn = 1<<63 - 1
		for i := int64(0); i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
			if a[i] < mn {
				mn = a[i]
			}
			if a[i] > mx {
				mx = a[i]
			}
		}
		if n > m {
			fmt.Fprintln(out, "NO")
			continue
		}
		need := n + sum + mx - mn
		if need <= m {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
