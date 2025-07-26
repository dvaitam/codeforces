package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// For each divisor k of n, it checks whether there exists an
// integer m >= 2 such that replacing every element a_i with
// a_i mod m results in all subarrays of length k being
// identical. The condition reduces to requiring a common
// divisor greater than 1 for the differences within each
// residue class modulo k.

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		divisors := make([]int, 0)
		for d := 1; d*d <= n; d++ {
			if n%d == 0 {
				divisors = append(divisors, d)
				if d*d != n {
					divisors = append(divisors, n/d)
				}
			}
		}

		ans := 0
		for _, k := range divisors {
			gcdAll := 0
			for r := 0; r < k && gcdAll != 1; r++ {
				base := arr[r]
				g := 0
				for i := r + k; i < n; i += k {
					diff := arr[i] - base
					if diff < 0 {
						diff = -diff
					}
					g = gcd(g, diff)
					if g == 1 {
						break
					}
				}
				gcdAll = gcd(gcdAll, g)
			}
			if gcdAll == 0 || gcdAll > 1 {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
