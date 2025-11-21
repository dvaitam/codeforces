package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	dp1 := [2]int64{}
	dp2 := [2][2]int64{}
	var ans int64

	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		cur := x & 1

		oldDp1 := dp1
		oldDp2 := dp2

		dp1[cur]++
		if dp1[cur] >= mod {
			dp1[cur] -= mod
		}

		for p := 0; p < 2; p++ {
			add := oldDp1[p]
			if add != 0 {
				dp2[p][cur] += add
				if dp2[p][cur] >= mod {
					dp2[p][cur] -= mod
				}
			}
		}

		for p := 0; p < 2; p++ {
			for q := 0; q < 2; q++ {
				add := oldDp2[p][q]
				if add != 0 && cur == (p^q) {
					dp2[q][cur] += add
					if dp2[q][cur] >= mod {
						dp2[q][cur] -= mod
					}
					ans += add
					if ans >= mod {
						ans -= mod
					}
				}
			}
		}
	}

	fmt.Println(ans % mod)
}
