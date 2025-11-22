package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

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

	// Euler sieve for phi.
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

	// tri[m] — number of good arrays of length 3 whose maximum value is exactly m.
	tri := make([]int64, maxM+1)
	for d := 2; d <= maxM; d++ {
		ph := int64(phi[d])
		for m := d; m <= maxM; m += d {
			contrib := int64(m-d) * ph % mod * 2 % mod
			tri[m] = (tri[m] + contrib) % mod
		}
	}

	pref := make([]int64, maxM+1)
	for i := 1; i <= maxM; i++ {
		pref[i] = (pref[i-1] + tri[i]) % mod
	}

	for _, m := range ms {
		m64 := int64(m) % mod
		ans := (m64 + m64*int64(m-1)%mod + pref[m]) % mod
		fmt.Fprintln(out, ans)
	}
}
