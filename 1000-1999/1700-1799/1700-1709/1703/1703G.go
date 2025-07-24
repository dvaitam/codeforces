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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	const MaxHalve = 32
	const NegInf int64 = -1 << 60

	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		dp := make([]int64, MaxHalve+1)
		for i := 0; i <= MaxHalve; i++ {
			dp[i] = NegInf
		}
		dp[0] = 0

		for _, val := range a {
			ndp := make([]int64, MaxHalve+1)
			for i := 0; i <= MaxHalve; i++ {
				ndp[i] = NegInf
			}
			for j := 0; j <= MaxHalve; j++ {
				if dp[j] == NegInf {
					continue
				}
				gainGood := int64(0)
				if j < 60 { // to avoid shifting by >=64
					gainGood = val >> uint(j)
				}
				cur := dp[j] + gainGood - k
				if cur > ndp[j] {
					ndp[j] = cur
				}
				if j+1 <= MaxHalve {
					gainBad := int64(0)
					if j+1 < 60 {
						gainBad = val >> uint(j+1)
					}
					cur2 := dp[j] + gainBad
					if cur2 > ndp[j+1] {
						ndp[j+1] = cur2
					}
				}
			}
			dp = ndp
		}

		ans := NegInf
		for _, v := range dp {
			if v > ans {
				ans = v
			}
		}
		fmt.Fprintln(out, ans)
	}
}
