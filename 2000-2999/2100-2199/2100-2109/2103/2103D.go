package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		maxA := 0
		for i := range a {
			fmt.Fscan(in, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
		}
		finalLevel := maxA + 1 // treat -1 as this level

		levels := make([][]int, finalLevel+1)
		for idx, v := range a {
			lvl := finalLevel
			if v != -1 {
				lvl = v
			}
			levels[lvl] = append(levels[lvl], idx)
		}

		seq := make([]int, len(levels[finalLevel]))
		copy(seq, levels[finalLevel]) // current order of surviving elements
		// assign initial values to deepest level
		valMap := make(map[int]int, n) // index -> raw value
		for i, idx := range seq {
			valMap[idx] = i
		}
		curMin, curMax := 0, len(seq)-1

		// build sequence from deeper to shallower levels
		for lvl := finalLevel - 1; lvl >= 1; lvl-- {
			newElems := levels[lvl]
			if len(seq) > 1 && len(newElems) < len(seq)-1 {
				// Guaranteed not to happen for valid input, but keep safe.
				return
			}
			tmp := make([]int, 0, len(seq)+len(newElems))
			tmp = append(tmp, seq[0])
			ptr := 0
			for i := 1; i < len(seq); i++ {
				tmp = append(tmp, newElems[ptr])
				ptr++
				tmp = append(tmp, seq[i])
			}
			for ptr < len(newElems) {
				tmp = append(tmp, newElems[ptr])
				ptr++
			}

			if lvl%2 == 1 { // odd iteration keeps local minima => new elements should be larger
				val := curMax
				for _, idx := range tmp {
					if _, ok := valMap[idx]; ok {
						continue
					}
					val++
					valMap[idx] = val
				}
				curMax = val
			} else { // even iteration keeps local maxima => new elements should be smaller
				val := curMin
				for _, idx := range tmp {
					if _, ok := valMap[idx]; ok {
						continue
					}
					val--
					valMap[idx] = val
				}
				curMin = val
			}

			seq = tmp
		}

		// collect values in current order and compress to 1..n to form a permutation
		vals := make([]int, n)
		for i, idx := range seq {
			vals[i] = valMap[idx]
		}
		sortedVals := append([]int(nil), vals...)
		sort.Ints(sortedVals)
		rank := make(map[int]int, n)
		for i, v := range sortedVals {
			if _, ok := rank[v]; !ok {
				rank[v] = i + 1
			}
		}

		for i, v := range vals {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, rank[v])
		}
		fmt.Fprintln(out)
	}
}
