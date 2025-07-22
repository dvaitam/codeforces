package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k, p int
	if _, err := fmt.Fscan(in, &n, &k, &p); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const NEG = -1 << 60
	best := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		best[i] = make([]int, p)
		for j := 0; j < p; j++ {
			best[i][j] = NEG
		}
	}
	best[0][0] = 0

	prefix := 0
	ans := NEG
	for idx := 0; idx < n; idx++ {
		prefix = (prefix + a[idx]) % p
		maxT := k
		if idx+1 < maxT {
			maxT = idx + 1
		}
		for t := maxT; t >= 1; t-- {
			val := NEG
			for r := 0; r < p; r++ {
				prev := best[t-1][r]
				if prev == NEG {
					continue
				}
				gain := prefix - r
				if gain < 0 {
					gain += p
				}
				candidate := prev + gain
				if candidate > val {
					val = candidate
				}
			}
			if val > best[t][prefix] {
				best[t][prefix] = val
			}
			if idx == n-1 && t == k {
				ans = val
			}
		}
	}
	if ans == NEG {
		ans = 0
	}
	fmt.Println(ans)
}
