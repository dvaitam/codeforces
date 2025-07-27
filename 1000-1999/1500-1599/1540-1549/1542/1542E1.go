package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var mod int64
	if _, err := fmt.Fscan(in, &n, &mod); err != nil {
		return
	}
	maxD := n * (n - 1) / 2
	offset := maxD
	size := 2*maxD + 1
	dp0 := make([]int64, size)
	dp1 := make([]int64, size)
	dp0[offset] = 1

	for m := n; m >= 1; m-- {
		next0 := make([]int64, size)
		next1 := make([]int64, size)
		for d := -maxD; d <= maxD; d++ {
			idx := d + offset
			if idx < 0 || idx >= size {
				continue
			}
			v0 := dp0[idx]
			if v0 != 0 {
				next0[idx] = (next0[idx] + v0*int64(m)) % mod
				for k := 1; k < m; k++ {
					delta := -k
					idx2 := d + delta + offset
					if idx2 >= 0 && idx2 < size {
						next1[idx2] = (next1[idx2] + v0*int64(m-k)) % mod
					}
				}
			}
			v1 := dp1[idx]
			if v1 != 0 {
				for delta := -m + 1; delta <= m-1; delta++ {
					idx2 := d + delta + offset
					if idx2 >= 0 && idx2 < size {
						cnt := m - absInt(delta)
						next1[idx2] = (next1[idx2] + v1*int64(cnt)) % mod
					}
				}
			}
		}
		dp0, dp1 = next0, next1
	}

	var ans int64
	for d := 1; d <= maxD; d++ {
		ans = (ans + dp1[d+offset]) % mod
	}
	fmt.Println(ans % mod)
}
