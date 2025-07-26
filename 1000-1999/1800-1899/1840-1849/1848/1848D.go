package main

import (
	"bufio"
	"fmt"
	"os"
)

func prefixSum(start, r int) int {
	cycle := []int{2, 4, 8, 6}
	// find index of start in cycle
	idx := 0
	for i, v := range cycle {
		if v == start {
			idx = i
			break
		}
	}
	sum := 0
	for i := 0; i < r; i++ {
		sum += cycle[(idx+i)%4]
	}
	return sum
}

func bestAfterCycle(kRem int64, sPre int64, startDigit int) int64 {
	if kRem <= 0 {
		return 0
	}
	var best int64
	maxR := int64(3)
	if kRem < 3 {
		maxR = kRem
	}
	for r := int64(0); r <= maxR; r++ {
		A := kRem - r
		B := sPre + int64(prefixSum(startDigit, int(r)))
		cMax := A / 4
		if cMax < 0 {
			cMax = 0
		}
		// candidate around optimal
		cOpt := (5*A - B) / 40
		candidates := []int64{0, cMax, cOpt - 2, cOpt - 1, cOpt, cOpt + 1, cOpt + 2}
		for _, c := range candidates {
			if c < 0 || c > cMax {
				continue
			}
			bFinal := B + 20*c
			mUsed := r + 4*c
			if mUsed > kRem {
				continue
			}
			disc := (kRem - mUsed) * bFinal
			if disc > best {
				best = disc
			}
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s, k int64
		fmt.Fscan(in, &s, &k)
		ans := s * k
		d := s % 10
		if d == 0 {
			fmt.Fprintln(out, ans)
			continue
		}
		if d == 5 {
			if k >= 1 {
				val := (s + 5) * (k - 1)
				if val > ans {
					ans = val
				}
			}
			fmt.Fprintln(out, ans)
			continue
		}
		// digits other than 0 or 5
		// case with no additional accumulation already covered by ans
		// option with prefix step
		if k >= 1 {
			sPre := s + d
			startDigit := int((d * 2) % 10)
			val := bestAfterCycle(k-1, sPre, startDigit)
			if val > ans {
				ans = val
			}
		}
		// option with immediate cycle (if digit already even)
		startDigit := int(d)
		if startDigit%2 == 0 {
			val := bestAfterCycle(k, s, startDigit)
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}
