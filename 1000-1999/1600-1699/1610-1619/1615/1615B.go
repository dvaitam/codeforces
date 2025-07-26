package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to problemB.txt for 1615B (Bitwise AND). The goal is to remove
// the minimum number of elements from the range [l, r] so that the bitwise
// AND of the remaining elements is non‑zero. For any bit position we can keep
// all numbers that have that bit set. We choose the bit that appears in the
// most numbers. The answer is the total count minus this maximum.

func countOnes(x, bit int) int {
	if x < 0 {
		return 0
	}
	cycle := 1 << (bit + 1)
	full := (x + 1) / cycle
	rem := (x + 1) % cycle
	ones := full * (1 << bit)
	if rem > (1 << bit) {
		ones += rem - (1 << bit)
	}
	return ones
}

func onesInRange(l, r, bit int) int {
	return countOnes(r, bit) - countOnes(l-1, bit)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		n := r - l + 1
		maxOnes := 0
		for bit := 0; bit < 20; bit++ {
			ones := onesInRange(l, r, bit)
			if ones > maxOnes {
				maxOnes = ones
			}
		}
		fmt.Fprintln(out, n-maxOnes)
	}
}
