package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([][4]int64, n)
	var sumA int64
	suitSums := make([]int64, 4)
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			fmt.Fscan(in, &a[i][j])
			sumA += a[i][j]
			suitSums[j] += a[i][j]
		}
	}

	b := make([]int64, 4)
	var sumB int64
	for j := 0; j < 4; j++ {
		fmt.Fscan(in, &b[j])
		sumB += b[j]
		suitSums[j] += b[j]
	}

	total := sumA + sumB
	T := total / int64(n)

	// prefix and suffix maxima of initial best suits
	maxSuit := make([]int64, n)
	for i := 0; i < n; i++ {
		mx := a[i][0]
		for j := 1; j < 4; j++ {
			if a[i][j] > mx {
				mx = a[i][j]
			}
		}
		maxSuit[i] = mx
	}
	prefix := make([]int64, n)
	if n > 0 {
		prefix[0] = maxSuit[0]
		for i := 1; i < n; i++ {
			if maxSuit[i] > prefix[i-1] {
				prefix[i] = maxSuit[i]
			} else {
				prefix[i] = prefix[i-1]
			}
		}
	}
	suffix := make([]int64, n)
	if n > 0 {
		suffix[n-1] = maxSuit[n-1]
		for i := n - 2; i >= 0; i-- {
			if maxSuit[i] > suffix[i+1] {
				suffix[i] = maxSuit[i]
			} else {
				suffix[i] = suffix[i+1]
			}
		}
	}

	results := make([]int64, n)

	for i := 0; i < n; i++ {
		baseOthers := int64(0)
		if i > 0 {
			baseOthers = prefix[i-1]
		}
		if i+1 < n {
			if suffix[i+1] > baseOthers {
				baseOthers = suffix[i+1]
			}
		}

		d := T
		for j := 0; j < 4; j++ {
			d -= a[i][j]
		}

		// S_other holds total cards of each suit excluding player i's current cards
		S_other := make([]int64, 4)
		for j := 0; j < 4; j++ {
			S_other[j] = suitSums[j] - a[i][j]
		}

		bestDiff := int64(0)

		for j0 := 0; j0 < 4; j0++ {
			// Copy arrays to modify
			Btemp := make([]int64, 4)
			copy(Btemp, b)
			Stemp := make([]int64, 4)
			copy(Stemp, S_other)
			x := make([]int64, 4)

			take := Btemp[j0]
			if take > d {
				take = d
			}
			x[j0] = take
			Btemp[j0] -= take
			Stemp[j0] -= take
			remaining := d - take

			for remaining > 0 {
				best := -1
				for s := 0; s < 4; s++ {
					if Btemp[s] > 0 {
						if best == -1 || Stemp[s] > Stemp[best] {
							best = s
						}
					}
				}
				if best == -1 {
					break
				}
				give := Btemp[best]
				if give > remaining {
					give = remaining
				}
				x[best] += give
				Btemp[best] -= give
				Stemp[best] -= give
				remaining -= give
			}

			// Player's best suit after allocation
			curMax := a[i][0] + x[0]
			for s := 1; s < 4; s++ {
				val := a[i][s] + x[s]
				if val > curMax {
					curMax = val
				}
			}

			y := baseOthers
			for s := 0; s < 4; s++ {
				val := (Stemp[s] + int64(n-1) - 1) / int64(n-1)
				if val > y {
					y = val
				}
			}

			diff := curMax - y
			if diff > bestDiff {
				bestDiff = diff
			}
		}

		if bestDiff < 0 {
			bestDiff = 0
		}
		results[i] = bestDiff
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, results[i])
	}
	fmt.Fprintln(out)
}
