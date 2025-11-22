package main

import (
	"bufio"
	"fmt"
	"os"
)

// Parity insight:
// A sequence is "good" iff positive counts are nondecreasing. The median depends
// only on the multiset of values. In the XOR, a multiset contributes once iff
// the number of its permutations is odd. The multinomial coefficient is odd
// exactly when the counts form a carry‑free partition of n in binary, meaning
// every bit of n is assigned to exactly one value (no overlapping 1-bits).
//
// Arrange the 1-bits of n from least to most significant. To keep counts
// nondecreasing, these bits can only be grouped into consecutive segments in
// that order; sums of segments are strictly increasing because each power of
// two exceeds the sum of all smaller powers. Positions with zero count can be
// inserted anywhere between segments.
//
// The median group must contain the most significant bit, hence it is always
// the last segment. Let its starting bit index be t. The number of ways to
// partition the first t bits into earlier segments equals 2^{t-1}, whose parity
// is 1 only for t = 0 or 1. So only two cases matter:
//   - Single segment containing all bits (always valid).
//   - Two segments where the first contains only the least significant set bit;
//     valid iff n has at least two set bits (i.e., n is not a power of two).
//
// Placement of g segments into m values is an increasing choice of positions.
// For g = 1, each of the m positions yields median = position (parity 1).
// For g = 2, with median at position v > 0, the first segment can be placed in
// any of v earlier positions; parity equals v mod 2. Therefore medians with odd
// v contribute.
//
// Final answer:
//   XOR_all_positions = XOR of 0..m-1
//   If n has at least two set bits, additionally XOR all odd numbers in [1, m-1].

func xorTo(x int64) int64 {
	switch x & 3 {
	case 0:
		return x
	case 1:
		return 1
	case 2:
		return x + 1
	default:
		return 0
	}
}

// XOR of odd numbers 1,3,...,t where t is odd.
func xorOddsUpTo(limit int64) int64 {
	if limit < 1 {
		return 0
	}
	t := limit
	if t%2 == 0 {
		t--
	}
	q := t / 2 // highest even index; evens count = q
	return xorTo(t) ^ (2 * xorTo(q))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var k int
		var m int64
		fmt.Fscan(in, &k, &m)
		var s string
		fmt.Fscan(in, &s)

		ones := 0
		for _, c := range s {
			if c == '1' {
				ones++
			}
		}

		ans := xorTo(m - 1) // case with single segment always contributes
		if ones >= 2 {
			ans ^= xorOddsUpTo(m - 1)
		}

		fmt.Fprintln(out, ans)
	}
}
