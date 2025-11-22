package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Greedy feasibility check explained in code comments below.
// We work with maximal alternating segments of the string. For a segment of
// alternating characters with even length 2k and different ends (pattern ABAB…
// or BABA…), the default tiling uses k pairs of the starting orientation. If
// we switch the segment to contain any pair of the opposite orientation, the
// total number of pairs drops by exactly 1 and the opposite orientation count
// becomes k-1 or k-2 (depending on whether we keep one original-orientation
// pair). This behavior lets us model a "conversion" of such a segment as
// losing one pair while reducing the dominant orientation usage by k-1 or k.
//
// Odd-length alternating segments contribute floor(len/2) totally flexible
// pairs that can be oriented as we like, limited only by the remaining
// capacities for AB and BA pairs.
//
// Algorithm outline for each test case:
// 1) Collect lengths k for even alternating segments starting with 'A'
//    (type AB) and starting with 'B' (type BA), and count flexible pairs from
//    odd segments.
// 2) Start with AB usage = sum(kAB), BA usage = sum(kBA), total pairs the same.
// 3) If AB usage exceeds the cap, greedily convert AB-segments (largest k
//    first) while BA capacity permits. Each conversion loses one pair and
//    reduces AB usage by k-1 or k (preferring the larger reduction only when
//    still above the cap). If excess remains, drop the extra AB pairs (each
//    drop loses one pair).
// 4) Symmetrically handle BA over-cap with BA-type segments.
// 5) Use remaining capacity to add flexible pairs from odd segments.
// 6) Check if the resulting maximum total pairs reaches the minimum number of
//    pairs required to respect single-letter limits: needPairs = max(A-a, B-b).
//    If yes, the split is possible.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		var aCap, bCap, abCap, baCap int
		fmt.Fscan(in, &aCap, &bCap, &abCap, &baCap)

		cntA, cntB := 0, 0
		for _, ch := range s {
			if ch == 'A' {
				cntA++
			} else {
				cntB++
			}
		}

		needPairs := cntA - aCap
		if tmp := cntB - bCap; tmp > needPairs {
			needPairs = tmp
		}
		if needPairs < 0 {
			needPairs = 0
		}

		kAB := make([]int, 0)
		kBA := make([]int, 0)
		flexPairs := 0

		// Split into maximal alternating segments
		for i := 0; i < len(s); {
			j := i
			for j+1 < len(s) && s[j+1] != s[j] {
				j++
			}
			length := j - i + 1
			if length%2 == 0 {
				// even alternating segment, ends differ
				k := length / 2
				if s[i] == 'A' {
					kAB = append(kAB, k)
				} else {
					kBA = append(kBA, k)
				}
			} else {
				flexPairs += length / 2
			}
			i = j + 1
		}

		abUse, baUse := 0, 0
		for _, v := range kAB {
			abUse += v
		}
		for _, v := range kBA {
			baUse += v
		}
		totalPairs := abUse + baUse

		sort.Slice(kAB, func(i, j int) bool { return kAB[i] > kAB[j] })
		sort.Slice(kBA, func(i, j int) bool { return kBA[i] > kBA[j] })

		// Reduce AB usage if necessary
		for _, k := range kAB {
			if abUse <= abCap {
				break
			}
			// Try full reduction (AB -> 0, BA -> k-1) if we still exceed cap after losing k-1
			if abUse-k >= abCap && baCap-baUse >= k-1 {
				abUse -= k
				baUse += k - 1
				totalPairs--
			} else if baCap-baUse >= k-2 {
				// Keep one AB, rest BA (needs k-2 BA capacity)
				abUse -= k - 1
				baUse += k - 2
				totalPairs--
			}
		}
		if abUse > abCap {
			totalPairs -= abUse - abCap
			abUse = abCap
		}

		// Reduce BA usage if necessary (symmetric)
		for _, k := range kBA {
			if baUse <= baCap {
				break
			}
			if baUse-k >= baCap && abCap-abUse >= k-1 {
				baUse -= k
				abUse += k - 1
				totalPairs--
			} else if abCap-abUse >= k-2 {
				baUse -= k - 1
				abUse += k - 2
				totalPairs--
			}
		}
		if baUse > baCap {
			totalPairs -= baUse - baCap
			baUse = baCap
		}

		remAB := abCap - abUse
		remBA := baCap - baUse
		addFlex := flexPairs
		if addFlex > remAB+remBA {
			addFlex = remAB + remBA
		}
		totalPairs += addFlex

		if totalPairs >= needPairs {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
