package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	ms := make([]int, t)
	maxM := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &ms[i])
		if ms[i] > maxM {
			maxM = ms[i]
		}
	}

	if maxM == 0 {
		for i := 0; i < t; i++ {
			fmt.Fprintln(out, 0)
		}
		return
	}

	phi := make([]int, maxM+1)
	for i := 0; i <= maxM; i++ {
		phi[i] = i
	}
	for i := 2; i <= maxM; i++ {
		if phi[i] == i {
			for j := i; j <= maxM; j += i {
				phi[j] -= phi[j] / i
			}
		}
	}

	tri := make([]int64, maxM+1)
	for mp := 2; mp <= maxM; mp++ {
		ph := int64(phi[mp])
		for m := mp; m <= maxM; m += mp {
			contrib := (int64(m-mp) % MOD) * ph % MOD * 2 % MOD
			tri[m] = (tri[m] + contrib) % MOD
		}
	}

	pref := make([]int64, maxM+1)
	for i := 1; i <= maxM; i++ {
		pref[i] = (pref[i-1] + tri[i]) % MOD
	}

	for _, m := range ms {
		m64 := int64(m)
		ans := (m64%MOD + (m64%MOD)*int64(m-1)%MOD + pref[m]) % MOD
		fmt.Fprintln(out, ans)
	}
}
