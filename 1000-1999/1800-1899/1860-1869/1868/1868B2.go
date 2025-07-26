package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const maxBits = 35

func solve(a []int64) bool {
	n := int64(len(a))
	var sum int64
	for _, v := range a {
		sum += v
	}
	if sum%n != 0 {
		return false
	}
	target := sum / n
	inc := make([]int64, maxBits)
	out := make([]int64, maxBits)
	posVar := make([]int64, maxBits)
	negVar := make([]int64, maxBits)

	for _, ai := range a {
		diff := target - ai
		if diff == 0 {
			continue
		}
		if diff > 0 {
			if diff&(diff-1) == 0 {
				k := bits.TrailingZeros64(uint64(diff))
				inc[k]++
				posVar[k]++
			} else {
				s := bits.TrailingZeros64(uint64(diff))
				val := diff + (1 << s)
				if val&(val-1) != 0 {
					return false
				}
				r := bits.Len64(uint64(val)) - 1
				if int64(1)<<uint(s) > ai+int64(1)<<uint(r) {
					return false
				}
				inc[r]++
				out[s]++
			}
		} else {
			y := -diff
			if y&(y-1) == 0 {
				k := bits.TrailingZeros64(uint64(y))
				out[k]++
				if ai >= int64(1)<<uint(k) {
					negVar[k]++
				}
			} else {
				r := bits.TrailingZeros64(uint64(y))
				val := y + (1 << r)
				if val&(val-1) != 0 {
					return false
				}
				s := bits.Len64(uint64(val)) - 1
				if int64(1)<<uint(s) > ai+int64(1)<<uint(r) {
					return false
				}
				inc[r]++
				out[s]++
			}
		}
	}

	diffCnt := make([]int64, maxBits)
	for i := 0; i < maxBits; i++ {
		diffCnt[i] = inc[i] - out[i]
	}

	for k := 0; k < maxBits-1; k++ {
		d := diffCnt[k]
		if d%2 != 0 {
			return false
		}
		need := d / 2
		yLow := int64(0)
		if -need > 0 {
			yLow = -need
		}
		yHigh := negVar[k]
		if posVar[k]-need < yHigh {
			yHigh = posVar[k] - need
		}
		if yLow > yHigh {
			return false
		}
		y := yLow
		x := need + y
		posVar[k] -= x
		negVar[k] -= y
		diffCnt[k+1] += x - y
	}
	if diffCnt[maxBits-1] != 0 {
		return false
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		if solve(arr) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
