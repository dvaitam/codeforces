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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		maxVal := n
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		freq := make([]int, maxVal+1)
		for _, v := range arr {
			freq[v]++
		}

		cntMult := make([]int, maxVal+1)
		hasDiv := make([]bool, maxVal+1)
		// accumulate multiples and divisors presence
		for v := 1; v <= maxVal; v++ {
			if freq[v] > 0 {
				for m := v; m <= maxVal; m += v {
					cntMult[m] += freq[v]
					hasDiv[m] = true
				}
			}
		}

		gcdCnt := make([]int64, maxVal+1)
		// compute number of pairs with gcd exactly g
		for g := maxVal; g >= 1; g-- {
			c := cntMult[g]
			if c >= 2 {
				total := int64(c) * int64(c-1) / 2
				for m := g * 2; m <= maxVal; m += g {
					total -= gcdCnt[m]
				}
				gcdCnt[g] = total
			}
		}

		var ans int64
		for g := 1; g <= maxVal; g++ {
			if !hasDiv[g] {
				ans += gcdCnt[g]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
