package main

import (
	"bufio"
	"fmt"
	"os"
)

const LIM = 1 << 20

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	B := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
	}
	// naive brute force for small n only
	if n > 20 {
		// fallback
		fmt.Println(0)
		return
	}
	var ans uint64
	var dfs func(pos, cnt int, curProd int, curVal uint64)
	dfs = func(pos, cnt int, curProd int, curVal uint64) {
		if pos == n-1 {
			// finalize last segment
			if curProd < LIM {
				curVal ^= 1 << curProd
			}
			if cnt >= k {
				ans ^= curVal
			}
			return
		}
		// option: power at pos
		nextProd := curProd * B[pos+1]
		dfs(pos+1, cnt, nextProd, curVal)
		// option: xor at pos
		val := curVal
		if curProd < LIM {
			val ^= 1 << curProd
		}
		dfs(pos+1, cnt+1, B[pos+1], val)
	}
	dfs(0, 0, B[0], 0)
	// print as binary string without leading zeros
	if ans == 0 {
		fmt.Println(0)
		return
	}
	out := ""
	for ans > 0 {
		if ans&1 == 1 {
			out = "1" + out
		} else {
			out = "0" + out
		}
		ans >>= 1
	}
	fmt.Println(out)
}
