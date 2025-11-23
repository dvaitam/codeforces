package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Set up fast I/O
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	// Read number of test cases
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		solve(reader, writer)
	}
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	var nonZeros []int // Stores 1-based indices of non-zero elements

	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
		if a[i] != 0 {
			nonZeros = append(nonZeros, i+1)
		}
	}

	// If the number of non-zero elements is odd, it's impossible to get a sum of 0.
	// Each non-zero element contributes ±1 to the alternating sum.
	// The parity of the total sum is equal to the parity of the count of non-zero elements.
	// 0 is even, so count must be even.
	if len(nonZeros)%2 != 0 {
		fmt.Fprintln(w, -1)
		return
	}

	var segments [][2]int
	segments = make([][2]int, 0, n)
	
	cur := 1
	
	// Process non-zero elements in pairs
	for k := 0; k < len(nonZeros); k += 2 {
		u := nonZeros[k]
		v := nonZeros[k+1]

		// Fill zeros before the current pair with single-element segments
		for j := cur; j < u; j++ {
			segments = append(segments, [2]int{j, j})
		}

		valU := a[u-1]
		valV := a[v-1]

		if valU + valV == 0 {
			// Case A: values are (1, -1) or (-1, 1).
			// We want their contributions to be +valU and +valV (which sum to 0).
			// To get +val, we start a segment at that index.
			// We isolate u and v in their own segments.
			segments = append(segments, [2]int{u, u})
			// Handle any zeros between u and v
			for j := u + 1; j < v; j++ {
				segments = append(segments, [2]int{j, j})
			}
			segments = append(segments, [2]int{v, v})
		} else {
			// Case B: values are (1, 1) or (-1, -1). Sum is not 0.
			// We need contributions +valU and -valV (so valU - valV = 0).
			// u is the start of a segment (contribution +).
			// v needs to be in a position such that its sign is flipped (-).
			// The sign of element at index x in segment [l, r] is (-1)^(x-l).
			// We need (-1)^(v-l) = -1, so v-l must be odd.
			
			if (v-u)%2 != 0 {
				// If v and u have different parity, v-u is odd.
				// A single segment [u, v] works. Start is u.
				// Sign of v is (-1)^(v-u) = -1.
				segments = append(segments, [2]int{u, v})
			} else {
				// If v and u have same parity, v-u is even.
				// We cannot use [u, v] because v would have sign +.
				// We split into [u, u] and [u+1, v].
				// [u, u] gives +valU.
				// [u+1, v] starts at u+1.
				// Since u, v same parity, u+1 and v have different parity.
				// So v-(u+1) is odd. Sign of v is -1.
				// Note: since u and v are consecutive non-zeros, u+1 is guaranteed to be a zero (or v itself, but parity check prevents v=u+1).
				// Actually, if v=u+1, diff is 1 (odd), handled above. So here v >= u+2, so u+1 is a zero.
				segments = append(segments, [2]int{u, u})
				segments = append(segments, [2]int{u + 1, v})
			}
		}
		cur = v + 1
	}

	// Fill remaining zeros after the last pair
	for j := cur; j <= n; j++ {
		segments = append(segments, [2]int{j, j})
	}

	// Output result
	fmt.Fprintln(w, len(segments))
	for _, seg := range segments {
		fmt.Fprintln(w, seg[0], seg[1])
	}
}