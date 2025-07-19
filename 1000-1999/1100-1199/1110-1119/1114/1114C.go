package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// calc returns the exponent of prime k in n! (factorial of n)
func calc(n, k int64) int64 {
	var res int64
	for n > 0 {
		res += n / k
		n /= k
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, b int64
	if _, err := fmt.Fscan(reader, &n, &b); err != nil {
		return
	}
	// copy of b for factorization
	bb := b
	ans := int64(math.MaxInt64)
	for i := int64(2); i*i <= bb; i++ {
		if bb%i == 0 {
			var cnt int64
			for bb%i == 0 {
				cnt++
				bb /= i
			}
			// compute how many times i divides n! and divide by its count in b
			val := calc(n, i) / cnt
			if val < ans {
				ans = val
			}
		}
	}
	// if remainder is a prime > 1
	if bb > 1 {
		val := calc(n, bb)
		if val < ans {
			ans = val
		}
	}
	fmt.Println(ans)
}
