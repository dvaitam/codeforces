package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// Evaluate P(2) using big integers
	p := big.NewInt(0)
	for i := n; i >= 0; i-- {
		p.Lsh(p, 1)
		p.Add(p, big.NewInt(a[i]))
	}

	y := new(big.Int).Set(p) // current value of P(2) >> j
	divisible := true        // whether 2^j divides P(2)
	ans := 0

	for j := 0; j <= n; j++ {
		if divisible && y.BitLen() <= 62 {
			val := y.Int64()
			diff := a[j] - val
			if diff >= -k && diff <= k {
				if j != n || diff != 0 {
					ans++
				}
			}
		}
		// prepare for next shift
		divisible = divisible && y.Bit(0) == 0
		y.Rsh(y, 1)
	}

	fmt.Println(ans)
}
