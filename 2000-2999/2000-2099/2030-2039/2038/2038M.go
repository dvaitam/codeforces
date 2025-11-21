package main

import (
	"bufio"
	"fmt"
	"os"
)

type State struct {
	rem     [4]uint8
	hand    [4]uint8
	remBad  uint8
	handBad uint8
}

var (
	n    int
	comb [][]float64
	memo map[State]float64
	inf  = 1e100
)

func prepareComb(maxN int) {
	comb = make([][]float64, maxN+1)
	for i := 0; i <= maxN; i++ {
		comb[i] = make([]float64, 6) // we only need up to choose 5
		comb[i][0] = 1
		for k := 1; k <= 5 && k <= i; k++ {
			comb[i][k] = comb[i-1][k-1] + comb[i-1][k]
		}
	}
}

func canon(st State) State {
	// sort suits to canonical order (descending rem then hand)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			ri, hi := st.rem[i], st.hand[i]
			rj, hj := st.rem[j], st.hand[j]
			if ri < rj || (ri == rj && hi < hj) {
				st.rem[i], st.rem[j] = rj, ri
				st.hand[i], st.hand[j] = hj, hi
			}
		}
	}
	// suits beyond n already zero, ok
	return st
}

func solve(st State) float64 {
	st = canon(st)
	if val, ok := memo[st]; ok {
		return val
	}

	// win check
	for i := 0; i < n; i++ {
		if st.hand[i] == 5 {
			memo[st] = 0
			return 0
		}
	}

	// determine possible suits (no discarded royals yet)
	possibleIdx := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if int(st.rem[i])+int(st.hand[i]) == 5 {
			possibleIdx = append(possibleIdx, i)
		}
	}
	if len(possibleIdx) == 0 {
		memo[st] = inf
		return inf
	}

	best := inf

	// Preprocess impossible suits: we'll set them to dropped automatically in each action
	// but since they can't help, we can add their remaining royals to remBad once here per action.

	// iterate over non-empty subsets of possible suits to keep
	subsetCount := 1 << len(possibleIdx)
	for mask := 1; mask < subsetCount; mask++ {
		var ns State
		ns.remBad = st.remBad
		// drop bad cards in hand unconditionally
		keptCards := 0

		// first handle all suits: copy values
		for i := 0; i < n; i++ {
			ns.rem[i] = st.rem[i]
			ns.hand[i] = st.hand[i]
		}

		// discard impossible suits (rem+hand<5)
		for i := 0; i < n; i++ {
			if int(st.rem[i])+int(st.hand[i]) < 5 {
				ns.remBad += ns.rem[i]
				ns.rem[i] = 0
				ns.hand[i] = 0
			}
		}

		// now process possible suits according to mask
		keep := make([]bool, n)
		for bit, idx := range possibleIdx {
			if mask&(1<<bit) != 0 {
				keep[idx] = true
			} else {
				// drop this suit
				ns.remBad += ns.rem[idx]
				ns.rem[idx] = 0
				ns.hand[idx] = 0
			}
		}

		// after drops compute keptCards
		for i := 0; i < n; i++ {
			keptCards += int(ns.hand[i])
		}

		if keptCards >= 5 {
			// no room to draw, would loop without progress
			continue
		}

		totalRem := int(ns.remBad)
		for i := 0; i < n; i++ {
			totalRem += int(ns.rem[i])
		}
		if totalRem == 0 {
			continue
		}

		need := 5 - keptCards
		if need > totalRem {
			need = totalRem
		}

		denom := comb[totalRem][need]
		exp := 0.0

		// collect kept suits indices for draw enumeration
		keptIdx := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if ns.rem[i] > 0 || keep[i] { // keep true implies suit still tracked even if rem 0
				if keep[i] {
					keptIdx = append(keptIdx, i)
				}
			}
		}
		// draw enumeration over kept suits + bad
		drawn := make([]int, len(keptIdx))
		var dfsDraw func(pos, left int)
		dfsDraw = func(pos, left int) {
			if pos == len(keptIdx) {
				// assign bad
				xBad := left
				if xBad < 0 || xBad > int(ns.remBad) {
					return
				}
				// compute probability weight numerator
				num := comb[int(ns.remBad)][xBad]
				for i := 0; i < len(keptIdx); i++ {
					idx := keptIdx[i]
					num *= comb[int(ns.rem[idx])][drawn[i]]
				}
				prob := num / denom

				// build next state
				var next State
				for i := 0; i < n; i++ {
					next.rem[i] = ns.rem[i]
					next.hand[i] = ns.hand[i]
				}
				next.remBad = ns.remBad - uint8(xBad)
				next.handBad = uint8(xBad)
				for i := 0; i < len(keptIdx); i++ {
					idx := keptIdx[i]
					next.rem[idx] -= uint8(drawn[i])
					next.hand[idx] += uint8(drawn[i])
				}

				val := solve(next)
				exp += prob * val
				return
			}
			idx := keptIdx[pos]
			maxTake := left
			if m := int(ns.rem[idx]); m < maxTake {
				maxTake = m
			}
			for x := 0; x <= maxTake; x++ {
				drawn[pos] = x
				dfsDraw(pos+1, left-x)
			}
		}
		dfsDraw(0, need)

		cand := 1 + exp
		if cand < best {
			best = cand
		}
	}

	memo[st] = best
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	_, err := fmt.Fscan(in, &n)
	if err != nil {
		return
	}

	prepareComb(52)
	memo = make(map[State]float64)

	// initial deck
	totalBad := 8 * n
	totalCards := 13 * n

	// enumerate initial 5-card hand distributions
	denom := comb[totalCards][5]
	ans := 0.0

	// counts per suit of royals drawn x[i], and bad xBad, sum 5
	var x [4]int
	var dfs func(pos, left int, remBad int, probNum float64)
	dfs = func(pos, left int, remBad int, probNum float64) {
		if pos == n {
			xBad := left
			if xBad < 0 || xBad > remBad {
				return
			}
			num := probNum * comb[remBad][xBad]

			var st State
			for i := 0; i < 4; i++ {
				if i < n {
					st.hand[i] = uint8(x[i])
					st.rem[i] = uint8(5 - x[i])
				} else {
					st.hand[i] = 0
					st.rem[i] = 0
				}
			}
			st.handBad = uint8(xBad)
			st.remBad = uint8(totalBad - xBad)

			prob := num / denom
			ans += prob * solve(st)
			return
		}
		maxTake := left
		if maxTake > 5 {
			maxTake = 5
		}
		if maxTake > 5 { // redundant
			maxTake = 5
		}
		if maxTake > 5-pos { // placeholder
		}
		for t := 0; t <= 5 && t <= left && t <= 5; t++ {
			if t > 5 {
				break
			}
			if t > 5 {
				break
			}
			if t > 5 {
				break
			}
			if t > 5 {
				break
			}
			if t > 5-pos {
				// ignore; no real constraint
			}
			x[pos] = t
			dfs(pos+1, left-t, remBad, probNum*comb[5][t])
		}
	}
	dfs(0, 5, totalBad, 1)

	fmt.Printf("%.9f\n", ans)
}
