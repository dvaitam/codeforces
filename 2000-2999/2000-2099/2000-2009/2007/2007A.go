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

	isCoprime := make([][]bool, 1001)
	for i := range isCoprime {
		isCoprime[i] = make([]bool, 1001)
	}
	mod := uint64(1)
	for a := 1; a <= 1000; a++ {
		for b := a + 1; b <= 1000; b++ {
			if gcd(a, b) == 1 {
				isCoprime[a][b] = true
				isCoprime[b][a] = true
			}
		}
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		cur := uint64(0)
		for i := l; i <= r; i++ {
			for j := i + 1; j <= r; j++ {
				if !isCoprime[i][j] {
					continue
				}
				for k := j + 1; k <= r; k++ {
					if isCoprime[i][k] && isCoprime[j][k] {
						cur++
					}
				}
			}
		}
		fmt.Fprintln(out, cur%mod)
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
