package main

import (
	"fmt"
	"math/bits"
)

func main() {
	var t int
	if _, err := fmt.Scan(&t); err != nil {
		return
	}

	// Precomputed answers for n = 2^k - 1, up to k = 25
	arr := []uint64{
		0, 1, 1, 5, 1, 21, 1, 85, 73, 341,
		89, 1365, 1, 5461, 4681, 21845, 1, 87381,
		1, 349525, 299593, 1398101, 178481, 5592405, 1082401,
	}

	for ; t > 0; t-- {
		var n uint64
		fmt.Scan(&n)
		if n == 0 {
			fmt.Println(0)
			continue
		}

		// Find smallest k such that 2^k > n
		k := bits.Len64(n)
		nextPow := uint64(1) << uint(k)
		if n+1 == nextPow {
			// n == 2^k - 1
			idx := k - 1
			if int(idx) < len(arr) {
				fmt.Println(arr[idx])
			} else {
				fmt.Println(nextPow - 1)
			}
		} else {
			fmt.Println(nextPow - 1)
		}
	}
}
