package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const (
	MaxXor = 1 << 13
	Words  = MaxXor / 64
	Limit  = 130
)

type Bitset [Words]uint64

func applyXor(dst *Bitset, src *Bitset, val int) {
	for i := 0; i < Words; i++ {
		w := src[i]
		for w != 0 {
			b := bits.TrailingZeros64(w)
			idx := i*64 + b
			nidx := idx ^ val
			dst[nidx>>6] |= 1 << (uint(nidx) & 63)
			w &= w - 1
		}
	}
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		dp := make([]Bitset, n+1)
		dp[0][0] = 1
		freq := make([]int, Limit+2)
		for r := 1; r <= n; r++ {
			for i := 0; i < Words; i++ {
				dp[r][i] = dp[r-1][i]
			}
			for i := range freq {
				freq[i] = 0
			}
			mex := 0
			for l := r; l >= 1 && r-l+1 <= Limit; l-- {
				v := a[l-1]
				if v <= Limit {
					freq[v]++
				}
				for mex <= Limit && freq[mex] > 0 {
					mex++
				}
				if mex > Limit {
					break
				}
				applyXor(&dp[r], &dp[l-1], mex)
			}
		}
		ans := 0
		for x := MaxXor - 1; x >= 0; x-- {
			if (dp[n][x>>6]>>(uint(x)&63))&1 != 0 {
				ans = x
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
