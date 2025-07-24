package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

// phi computes Euler's totient of n.
func phi(n int64) int64 {
	result := n
	m := n
	for i := int64(2); i*i <= m; i++ {
		if m%i == 0 {
			for m%i == 0 {
				m /= i
			}
			result = result / i * (i - 1)
		}
	}
	if m > 1 {
		result = result / m * (m - 1)
	}
	return result
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	fmt.Fscan(in, &n, &k)
	// For any k >= 1, Fk(n) equals phi(n)
	ans := phi(n) % mod
	fmt.Println(ans)
}
