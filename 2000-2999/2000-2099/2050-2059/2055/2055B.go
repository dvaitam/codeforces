package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const infT = int64(1 << 60)

func alignParity(x int64, parity int64) int64 {
	if (x & 1) != parity {
		return x + 1
	}
	return x
}

func ceilDiv(x, y int64) int64 {
	if y <= 0 {
		panic("non-positive divisor")
	}
	if x >= 0 {
		return (x + y - 1) / y
	}
	return x / y
}

func canCraft(a, b []int64) bool {
	n := len(a)
	if n == 1 {
		return a[0] >= b[0]
	}
	sumA, sumB := int64(0), int64(0)
	need := int64(0)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		sumA += a[i]
		sumB += b[i]
		diff := b[i] - a[i]
		if diff > need {
			need = diff
		}
		c[i] = a[i] - b[i]
	}
	if need < 0 {
		need = 0
	}
	if sumA < sumB {
		return false
	}

	cSorted := append([]int64(nil), c...)
	sort.Slice(cSorted, func(i, j int) bool { return cSorted[i] < cSorted[j] })

	prefSum := make([]int64, n+1)
	prefOdd := make([]int64, n+1)
	for i, val := range cSorted {
		prefSum[i+1] = prefSum[i] + val
		prefOdd[i+1] = prefOdd[i]
		if val&1 != 0 {
			prefOdd[i+1]++
		}
	}

	cExt := append(append([]int64(nil), cSorted...), infT)

	for m := 0; m <= n; m++ {
		lower := need
		if m > 0 {
			if tmp := cSorted[m-1] + 1; tmp > lower {
				lower = tmp
			}
		}
		upper := cExt[m]
		if lower > upper {
			continue
		}
		sumC := prefSum[m]
		oddCnt := prefOdd[m]
		for parity := int64(0); parity < 2; parity++ {
			start := alignParity(lower, parity)
			if start > upper {
				continue
			}
			var mismatch int64
			if parity == 0 {
				mismatch = oddCnt
			} else {
				mismatch = int64(m) - oddCnt
			}
			if m < 2 {
				denom := int64(2 - m)
				bound := ceilDiv(mismatch-sumC, denom)
				if start < bound {
					start = alignParity(bound, parity)
				}
				if start <= upper {
					return true
				}
			} else if m == 2 {
				if sumC-mismatch >= 0 {
					return true
				}
			} else {
				numerator := sumC - mismatch
				if numerator < 0 {
					continue
				}
				limit := numerator / int64(m-2)
				end := upper
				if limit < end {
					end = limit
				}
				if start <= end {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		if canCraft(a, b) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
