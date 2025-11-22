package main

import (
	"bufio"
	"fmt"
	"os"
)

const modGCD = 998244353

func gcd(a, b int64) int64 {
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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	const scanLimit int64 = 200000

	for ; T > 0; T-- {
		var l, r, G int64
		fmt.Fscan(in, &l, &r, &G)

		L := (l + G - 1) / G // ceil
		R := r / G           // floor

		if L > R {
			fmt.Fprintln(out, "-1 -1")
			continue
		}

		if L == R {
			if L == 1 {
				val := L * G
				fmt.Fprintln(out, val, val)
			} else {
				fmt.Fprintln(out, "-1 -1")
			}
			continue
		}

		var bestA, bestB int64 = -1, -1
		var bestDiff int64 = -1

		maxOffset := R - L
		if maxOffset > scanLimit {
			maxOffset = scanLimit
		}

		// Try expanding from both ends with small offsets.
		for s := int64(0); s <= maxOffset; s++ {
			a := L + s
			if a <= R && gcd(a, R) == 1 {
				diff := R - a
				if diff > bestDiff || (diff == bestDiff && a < bestA) {
					bestDiff = diff
					bestA = a
					bestB = R
				}
			}

			b := R - s
			if b >= L && gcd(L, b) == 1 {
				diff := b - L
				if diff > bestDiff || (diff == bestDiff && L < bestA) {
					bestDiff = diff
					bestA = L
					bestB = b
				}
			}

			// If we already reached the maximal possible distance, we can stop.
			if bestDiff == R-L {
				break
			}
		}

		if bestDiff == -1 {
			// Fallback to consecutive numbers, which are guaranteed to be coprime.
			bestA = L
			bestB = L + 1
		}

		fmt.Fprintln(out, bestA*G, bestB*G)
	}
}
