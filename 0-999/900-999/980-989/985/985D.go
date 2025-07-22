package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxSum(len, H int64) int64 {
	if len <= H {
		return len * (len + 1) / 2
	}
	k := (len - H + 2) / 2
	return k*H + k*(k-1)/2 + (len-k)*(len-k+1)/2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, H int64
	fmt.Fscan(in, &n, &H)
	lo, hi := int64(1), int64(2000000000)
	for lo < hi {
		mid := (lo + hi) / 2
		if maxSum(mid, H) >= n {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Fprintln(out, lo)
}
