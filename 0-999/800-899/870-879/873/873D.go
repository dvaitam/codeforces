package main

import (
	"bufio"
	"fmt"
	"os"
)

var out *bufio.Writer

func solve(l, r, k, L, R int) {
	if k == 1 {
		for i := L; i <= R; i++ {
			fmt.Fprint(out, i, " ")
		}
		return
	}
	mid := (l + r - 1) >> 1
	Mid := (L + R + 2) >> 1
	leftCnt := mid - l + 1
	if 2*leftCnt-1 >= k-2 {
		solve(l, mid, k-2, Mid, R)
		solve(mid+1, r, 1, L, Mid-1)
	} else {
		solve(l, mid, 2*leftCnt-1, Mid, R)
		solve(mid+1, r, k-2*leftCnt, L, Mid-1)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k%2 == 0 || n*2 <= k {
		fmt.Fprintln(out, -1)
		return
	}
	solve(1, n, k, 1, n)
	fmt.Fprintln(out)
}
