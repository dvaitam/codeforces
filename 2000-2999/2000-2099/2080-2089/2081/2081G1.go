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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	phi := make([]uint32, n+1)
	phi[1] = 1
	primes := make([]int, 0)
	isComp := make([]bool, n+1)

	for i := 2; i <= n; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			phi[i] = uint32(i - 1)
		}
		for _, p := range primes {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				phi[i*p] = phi[i] * uint32(p)
				break
			}
			phi[i*p] = phi[i] * uint32(p-1)
		}
	}

	var ans uint32
	for i := 1; i <= n; i++ {
		mod := uint32(i % int(phi[i]))
		ans += mod
	}

	fmt.Fprintln(out, ans)
}
