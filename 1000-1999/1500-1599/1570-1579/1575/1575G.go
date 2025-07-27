package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	maxA := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	maxV := maxA
	if n > maxV {
		maxV = n
	}

	phi := make([]int, maxV+1)
	for i := 0; i <= maxV; i++ {
		phi[i] = i
	}
	for i := 2; i <= maxV; i++ {
		if phi[i] == i {
			for j := i; j <= maxV; j += i {
				phi[j] -= phi[j] / i
			}
		}
	}

	divs := make([][]int, maxA+1)
	for d := 1; d <= maxA; d++ {
		for m := d; m <= maxA; m += d {
			divs[m] = append(divs[m], d)
		}
	}

	cnt := make([]int, maxA+1)
	used := make([]int, 0, 128)
	var ans int64 = 0

	for g := 1; g <= n; g++ {
		used = used[:0]
		for idx := g; idx <= n; idx += g {
			v := a[idx]
			for _, d := range divs[v] {
				if cnt[d] == 0 {
					used = append(used, d)
				}
				cnt[d]++
			}
		}
		var S int64 = 0
		for _, d := range used {
			c := cnt[d]
			S = (S + int64(phi[d])*int64(c)*int64(c)) % MOD
		}
		ans = (ans + int64(phi[g])*S) % MOD
		for _, d := range used {
			cnt[d] = 0
		}
	}

	fmt.Println(ans % MOD)
}
