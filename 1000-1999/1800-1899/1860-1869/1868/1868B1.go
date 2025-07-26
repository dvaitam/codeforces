package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func solve() bool {
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return false
	}
	a := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum += a[i]
	}
	if sum%int64(n) != 0 {
		return false
	}
	avg := sum / int64(n)

	inCnt := make(map[int]int)
	outCnt := make(map[int]int)

	for _, v := range a {
		diff := avg - v
		if diff == 0 {
			inCnt[0]++
			outCnt[0]++
			continue
		}
		if diff > 0 {
			x := diff
			tz := bits.TrailingZeros64(uint64(x))
			y := x + (1 << tz)
			if y&(y-1) != 0 {
				return false
			}
			aExp := bits.TrailingZeros64(uint64(y))
			bExp := tz
			inCnt[int(aExp)]++
			outCnt[int(bExp)]++
		} else {
			x := -diff
			tz := bits.TrailingZeros64(uint64(x))
			y := x + (1 << tz)
			if y&(y-1) != 0 {
				return false
			}
			aExp := bits.TrailingZeros64(uint64(y))
			bExp := tz
			inCnt[int(bExp)]++
			outCnt[int(aExp)]++
		}
	}
	for k, v := range inCnt {
		if outCnt[k] != v {
			return false
		}
	}
	for k, v := range outCnt {
		if inCnt[k] != v {
			return false
		}
	}
	return true
}

func main() {
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		if solve() {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
