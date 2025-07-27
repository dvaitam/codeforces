package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		var l1, r1, l2, r2 int64
		fmt.Fscan(reader, &n, &k)
		fmt.Fscan(reader, &l1, &r1)
		fmt.Fscan(reader, &l2, &r2)

		// ensure l1 <= l2 for convenience
		if l1 > l2 {
			l1, l2 = l2, l1
			r1, r2 = r2, r1
		}

		intersection := max64(0, min64(r1, r2)-max64(l1, l2))
		unionLen := max64(r1, r2) - min64(l1, l2)
		gap := int64(0)
		if r1 < l2 {
			gap = l2 - r1
		} else if r2 < l1 {
			gap = l1 - r2
		}

		base := int64(n) * intersection
		if k <= base {
			fmt.Fprintln(writer, 0)
			continue
		}
		need := k - base

		extraWithin := unionLen - intersection
		if gap > 0 {
			extraWithin = unionLen
		}

		ans := int64(1<<62 - 1)
		for i := 1; i <= n; i++ {
			cost := int64(0)
			if gap > 0 {
				cost += gap * int64(i)
			}
			gainLimit := extraWithin * int64(i)
			take := min64(need, gainLimit)
			cost += take
			remain := need - take
			cost += 2 * remain
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
