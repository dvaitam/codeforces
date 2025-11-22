package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func ceilLog2(x int64) int64 {
	if x <= 1 {
		return 0
	}
	return int64(bits.Len64(uint64(x - 1)))
}

func g(h, w int64) int64 {
	return ceilLog2(h) + ceilLog2(w)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, a, b int64
		fmt.Fscan(in, &n, &m, &a, &b)

		// If already a single cell, zero turns (though constraints prevent this).
		if n == 1 && m == 1 {
			fmt.Fprintln(out, 0)
			continue
		}

		const inf int64 = 1 << 60
		ans := inf

		// Cut horizontally, keep the side containing row a.
		if n > 1 {
			h := a
			if n-a+1 < h {
				h = n - a + 1
			}
			val := 1 + g(h, m)
			if val < ans {
				ans = val
			}
		}

		// Cut vertically, keep the side containing column b.
		if m > 1 {
			w := b
			if m-b+1 < w {
				w = m - b + 1
			}
			val := 1 + g(n, w)
			if val < ans {
				ans = val
			}
		}

		fmt.Fprintln(out, ans)
	}
}
