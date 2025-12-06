package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// Using bufio for fast I/O
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	// Read number of test cases
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		solve(writer, n, k)
	}
}

func solve(w *bufio.Writer, n int, k int64) {
	// If only one number is allowed, we must output k itself.
	if n == 1 {
		fmt.Fprintln(w, k)
		return
	}

	// Calculate the number of set bits in k
	maxOnes := bits.OnesCount64(uint64(k))
	bestP := -1

	// We iterate over each bit position p. If bit p is set in k,
	// we consider the strategy of splitting the term 2^p into a sum that fills
	// all lower bits (0 to p-1).
	// The term 2^p is effectively replaced by (2^p - 1) + 1.
	// We assign (2^p - 1) to one number (a1), and add 1 to the rest of k (a2).
	// This guarantees bits 0 to p-1 are set in the OR result.
	// Bit p will be 0 in a1. In a2, bit p is likely 0 (unless carries propagate,
	// but that only happens if lower bits were already full, in which case this
	// split wouldn't be optimal).
	// The new popcount is approximately: (popcount of k above p) + p.
	
	// Since k <= 10^9, 30 bits are sufficient, but we check up to 62 for safety.
	for i := 0; i < 62; i++ {
		if (k>>i)&1 == 1 {
			// Calculate potential ones count if we sacrifice bit i to fill 0..i-1
			// Bits above i are preserved. Bit i is lost. Bits 0..i-1 become 1.
			currentOnes := bits.OnesCount64(uint64(k>>(i+1))) + i
			
			if currentOnes > maxOnes {
				maxOnes = currentOnes
				bestP = i
			}
		}
	}

	if bestP != -1 {
		// We found a split that increases the number of 1s.
		// a1 takes the lower mask 2^bestP - 1
		// a2 takes the remainder
		a1 := (int64(1) << bestP) - 1
		a2 := k - a1
		
		// Output a1, a2, and n-2 zeros
		fmt.Fprint(w, a1, " ", a2)
		for j := 0; j < n-2; j++ {
			fmt.Fprint(w, " 0")
		}
		fmt.Fprintln(w)
	} else {
		// No split improves the result, output k and n-1 zeros
		fmt.Fprint(w, k)
		for j := 0; j < n-1; j++ {
			fmt.Fprint(w, " 0")
		}
		fmt.Fprintln(w)
	}
}