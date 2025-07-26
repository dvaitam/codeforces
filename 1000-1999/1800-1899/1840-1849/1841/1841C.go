package main

import (
	"bufio"
	"fmt"
	"os"
)

var val = []int64{1, 10, 100, 1000, 10000}

const negInf int64 = -1 << 60

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solve(s string) int64 {
	n := len(s)
	dp0 := make([]int64, 5) // no change used
	dp1 := make([]int64, 5) // change already used
	for i := 0; i < 5; i++ {
		dp0[i] = 0
		dp1[i] = negInf
	}
	for idx := n - 1; idx >= 0; idx-- {
		orig := int(s[idx] - 'A')
		newdp0 := [5]int64{negInf, negInf, negInf, negInf, negInf}
		newdp1 := [5]int64{negInf, negInf, negInf, negInf, negInf}
		for pm := 0; pm < 5; pm++ {
			if dp0[pm] != negInf {
				sign := val[orig]
				if orig < pm {
					sign = -sign
				}
				nm := pm
				if orig > nm {
					nm = orig
				}
				newdp0[nm] = max64(newdp0[nm], sign+dp0[pm])
				for ni := 0; ni < 5; ni++ {
					sign2 := val[ni]
					if ni < pm {
						sign2 = -sign2
					}
					nm2 := pm
					if ni > nm2 {
						nm2 = ni
					}
					newdp1[nm2] = max64(newdp1[nm2], sign2+dp0[pm])
				}
			}
			if dp1[pm] != negInf {
				sign := val[orig]
				if orig < pm {
					sign = -sign
				}
				nm := pm
				if orig > nm {
					nm = orig
				}
				newdp1[nm] = max64(newdp1[nm], sign+dp1[pm])
			}
		}
		for j := 0; j < 5; j++ {
			dp0[j] = newdp0[j]
			dp1[j] = newdp1[j]
		}
	}
	ans := negInf
	for i := 0; i < 5; i++ {
		ans = max64(ans, dp0[i])
		ans = max64(ans, dp1[i])
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solve(s))
	}
}
