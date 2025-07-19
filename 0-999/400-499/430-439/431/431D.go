package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	dp [70][70]int64
	x  [70]int
	K  int
)

func dfs(pos, cnt int, flag bool) int64 {
	if pos == 0 {
		if cnt == K {
			return 1
		}
		return 0
	}
	if !flag && dp[pos][cnt] != -1 {
		return dp[pos][cnt]
	}
	// determine upper bound for this bit
	u := 1
	if flag {
		u = x[pos]
	}
	var ret int64
	for i := 0; i <= u; i++ {
		ret += dfs(pos-1, cnt+i, flag && i == u)
	}
	if !flag {
		dp[pos][cnt] = ret
	}
	return ret
}

func ju(mid int64) int64 {
	 e := 0
	// extract bits of mid into x[1..]
	for mid > 0 {
		e++
		x[e] = int(mid & 1)
		mid >>= 1
	}
	// count numbers with same bit length <= original mid
	tmp1 := dfs(e, 0, true)
	// count numbers with bit lengths less than original mid
	e++
	for i := e; i >= 2; i-- {
		x[i] = x[i-1]
	}
	x[1] = 0
	tmp2 := dfs(e, 0, true)
	return tmp2 - tmp1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		var m int64
		var k int
		if _, err := fmt.Fscan(reader, &m, &k); err != nil {
			break
		}
		// initialize dp cache
		for i := range dp {
			for j := range dp[i] {
				dp[i][j] = -1
			}
		}
		K = k
		l, r := int64(1), int64(1e18+1)
		for l < r {
			mid := (l + r) >> 1
			if ju(mid) >= m {
				r = mid
			} else {
				l = mid + 1
			}
		}
		fmt.Fprintln(writer, l)
	}
}
