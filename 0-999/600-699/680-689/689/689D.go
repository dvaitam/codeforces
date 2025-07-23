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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	// prepare log table
	lg := make([]int, n+1)
	for i := 2; i <= n; i++ {
		lg[i] = lg[i/2] + 1
	}
	K := lg[n] + 1

	stMax := make([][]int, K)
	stMin := make([][]int, K)
	for k := 0; k < K; k++ {
		stMax[k] = make([]int, n+1)
		stMin[k] = make([]int, n+1)
	}
	for i := 1; i <= n; i++ {
		stMax[0][i] = a[i]
		stMin[0][i] = b[i]
	}
	for k := 1; (1 << k) <= n; k++ {
		step := 1 << (k - 1)
		for i := 1; i+(1<<k)-1 <= n; i++ {
			x := stMax[k-1][i]
			y := stMax[k-1][i+step]
			if x >= y {
				stMax[k][i] = x
			} else {
				stMax[k][i] = y
			}
			x = stMin[k-1][i]
			y = stMin[k-1][i+step]
			if x <= y {
				stMin[k][i] = x
			} else {
				stMin[k][i] = y
			}
		}
	}

	getMax := func(l, r int) int {
		j := lg[r-l+1]
		x := stMax[j][l]
		y := stMax[j][r-(1<<j)+1]
		if x >= y {
			return x
		}
		return y
	}
	getMin := func(l, r int) int {
		j := lg[r-l+1]
		x := stMin[j][l]
		y := stMin[j][r-(1<<j)+1]
		if x <= y {
			return x
		}
		return y
	}

	var ans int64
	for l := 1; l <= n; l++ {
		// first r where max >= min
		lo, hi := l, n
		pos1 := n + 1
		for lo <= hi {
			mid := (lo + hi) >> 1
			if getMax(l, mid) >= getMin(l, mid) {
				pos1 = mid
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
		if pos1 > n || getMax(l, pos1) != getMin(l, pos1) {
			continue
		}
		// last r where max <= min (which here means equal)
		lo, hi = pos1, n
		pos2 := pos1
		for lo <= hi {
			mid := (lo + hi) >> 1
			if getMax(l, mid) <= getMin(l, mid) {
				pos2 = mid
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
		ans += int64(pos2 - pos1 + 1)
	}
	fmt.Fprintln(out, ans)
}
