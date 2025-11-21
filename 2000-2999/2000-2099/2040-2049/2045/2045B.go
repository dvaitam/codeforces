package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var N, D, S int64
	fmt.Fscan(in, &N, &D, &S)
	if S > N {
		S = N
	}
	if S > D {
		fmt.Println(S)
		return
	}

	maxG := min64(D, N)
	maxK := maxG / S
	ans := S
	var k int64 = 1
	for k <= maxK {
		g := S * k
		Dval := D / g
		Nval := N / g
		mult := Dval + 1
		if mult > Nval {
			mult = Nval
		}
		if mult < 1 {
			mult = 1
		}
		candidate := g * mult
		if candidate > ans {
			ans = candidate
		}
		boundD := D / (S * Dval)
		boundN := N / (S * Nval)
		nextK := boundD
		if boundN < nextK {
			nextK = boundN
		}
		if nextK > maxK {
			nextK = maxK
		}
		if nextK > k {
			bestK := nextK
			gBest := S * bestK
			candidate = gBest * mult
			if candidate > ans {
				ans = candidate
			}
		}
		k = nextK + 1
	}
	fmt.Println(ans)
}
