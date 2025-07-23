package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	// Check minimal possible sum for gcd = 1
	minSum := big.NewInt(0).Mul(big.NewInt(k), big.NewInt(k+1))
	minSum.Div(minSum, big.NewInt(2))
	if minSum.Cmp(big.NewInt(n)) > 0 {
		fmt.Println(-1)
		return
	}

	var best int64
	for d := int64(1); d*d <= n; d++ {
		if n%d == 0 {
			if n/d >= k*(k+1)/2 && d > best {
				best = d
			}
			other := n / d
			if d >= k*(k+1)/2 && other > best {
				best = other
			}
		}
	}

	if best == 0 {
		fmt.Println(-1)
		return
	}

	sum := int64(0)
	for i := int64(1); i < k; i++ {
		v := best * i
		fmt.Print(v, " ")
		sum += v
	}
	fmt.Println(n - sum)
}
