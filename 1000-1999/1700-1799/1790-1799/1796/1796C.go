package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// For each test case we find the maximum size of a beautiful set
// (numbers totally ordered by divisibility) within [l, r] and
// count how many such sets of maximum size exist.
//
// Observations:
// - If we arrange the elements of a beautiful set in ascending order,
//   every element divides all following ones. To maximize the length
//   we prefer multiplying by 2 at each step, as this yields the
//   smallest possible growth.
// - Let k be the maximum size. Then l * 2^(k-1) \le r < l * 2^k.
//   Any set of size k consists of a starting value s and k-1
//   multipliers greater than 1. Since r/l < 2^k, among these
//   multipliers at most one can be greater than 2, and this value can
//   only be 3. Therefore valid sets are of two types:
//   1) all multipliers equal 2;
//   2) exactly one multiplier equal 3 and the rest equal 2.
// - The number of starting values s for type (1) is
//     floor(r / 2^(k-1)) - l + 1.
//   For type (2) we need s * 3 * 2^(k-2) \le r and we can place the
//   multiplier 3 in any of the k-1 positions, hence
//     (k-1) * max(0, floor(r / (3*2^(k-2)))-l+1).
// The final answer is the sum of these counts modulo 998244353.

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(reader, &l, &r)

		// Compute maximum size k
		k := int64(0)
		v := l
		for v <= r {
			k++
			v *= 2
		}

		// Count sets consisting only of doubling steps
		pow2 := int64(1) << (k - 1)
		count1 := r/pow2 - l + 1
		if count1 < 0 {
			count1 = 0
		}

		// Count sets containing exactly one tripling step
		count2 := int64(0)
		if k >= 2 {
			pow3 := int64(3) * (int64(1) << (k - 2))
			c := r/pow3 - l + 1
			if c > 0 {
				count2 = c * (k - 1)
			}
		}

		total := (count1 + count2) % mod
		fmt.Fprintf(writer, "%d %d\n", k, total)
	}
}
