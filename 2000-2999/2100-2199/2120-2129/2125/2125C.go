package main

import (
	"bufio"
	"fmt"
	"os"
)

var primes = []int64{2, 3, 5, 7}

func countCoprime(n int64) int64 {
	if n <= 0 {
		return 0
	}
	var res int64
	for mask := 0; mask < 1<<len(primes); mask++ {
		mult := int64(1)
		sign := int64(1)
		for i := 0; i < len(primes); i++ {
			if mask&(1<<i) != 0 {
				mult *= primes[i]
				sign = -sign
			}
		}
		res += sign * (n / mult)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		ans := countCoprime(r) - countCoprime(l-1)
		fmt.Fprintln(out, ans)
	}
}
