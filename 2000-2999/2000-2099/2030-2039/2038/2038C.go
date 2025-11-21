package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type entry struct {
	val   int64
	pairs int
}

type axisPair struct {
	i, j int
	diff int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		freq := make(map[int64]int)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			freq[x]++
		}

		all := make([]entry, 0, len(freq))
		totalPairs := 0
		for v, c := range freq {
			p := c / 2
			if p > 0 {
				all = append(all, entry{val: v, pairs: p})
				totalPairs += p
			}
		}

		if totalPairs < 4 {
			fmt.Fprintln(out, "NO")
			continue
		}

		sort.Slice(all, func(i, j int) bool {
			if all[i].val == all[j].val {
				return all[i].pairs > all[j].pairs
			}
			return all[i].val < all[j].val
		})

		candidates := selectCandidates(all, 30)

		bestArea, bestXP, bestYP := findBestRectangle(candidates)
		if bestArea < 0 {
			fmt.Fprintln(out, "NO")
			continue
		}

		x1 := candidates[bestXP.i].val
		x2 := candidates[bestXP.j].val
		y1 := candidates[bestYP.i].val
		y2 := candidates[bestYP.j].val

		fmt.Fprintln(out, "YES")
		fmt.Fprintf(out, "%d %d %d %d %d %d %d %d\n", x1, y1, x1, y2, x2, y1, x2, y2)
	}
}

func selectCandidates(all []entry, limit int) []entry {
	if len(all) <= limit*2 {
		cpy := make([]entry, len(all))
		copy(cpy, all)
		sort.Slice(cpy, func(i, j int) bool {
			return cpy[i].val < cpy[j].val
		})
		return cpy
	}

	selected := make(map[int]struct{})
	for i := 0; i < limit && i < len(all); i++ {
		selected[i] = struct{}{}
	}
	for i := len(all) - limit; i < len(all); i++ {
		if i >= 0 {
			selected[i] = struct{}{}
		}
	}

	idx := make([]int, len(all))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		if all[idx[i]].pairs == all[idx[j]].pairs {
			return all[idx[i]].val < all[idx[j]].val
		}
		return all[idx[i]].pairs > all[idx[j]].pairs
	})
	for i := 0; i < limit && i < len(idx); i++ {
		selected[idx[i]] = struct{}{}
	}

	result := make([]entry, 0, len(selected))
	for i := range all {
		if _, ok := selected[i]; ok {
			result = append(result, all[i])
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].val < result[j].val
	})
	return result
}

func findBestRectangle(cand []entry) (int64, axisPair, axisPair) {
	m := len(cand)
	if m == 0 {
		return -1, axisPair{}, axisPair{}
	}

	values := make([]int64, m)
	avail := make([]int, m)
	totalAvail := 0
	for i := 0; i < m; i++ {
		values[i] = cand[i].val
		if cand[i].pairs > 4 {
			avail[i] = 4
		} else {
			avail[i] = cand[i].pairs
		}
		totalAvail += avail[i]
	}
	if totalAvail < 4 {
		return -1, axisPair{}, axisPair{}
	}

	pairs := make([]axisPair, 0, m*m/2)
	for i := 0; i < m; i++ {
		if avail[i] == 0 {
			continue
		}
		for j := i; j < m; j++ {
			if i == j {
				if avail[i] < 2 {
					continue
				}
			} else {
				if avail[j] == 0 {
					continue
				}
			}
			diff := values[j] - values[i]
			if diff < 0 {
				diff = -diff
			}
			pairs = append(pairs, axisPair{i: i, j: j, diff: diff})
		}
	}

	bestArea := int64(-1)
	var bestXP, bestYP axisPair

	for i := 0; i < len(pairs); i++ {
		for j := i; j < len(pairs); j++ {
			if validCombination(pairs[i], pairs[j], avail) {
				area := pairs[i].diff * pairs[j].diff
				if area > bestArea {
					bestArea = area
					bestXP = pairs[i]
					bestYP = pairs[j]
				}
			}
		}
	}

	return bestArea, bestXP, bestYP
}

func validCombination(a, b axisPair, avail []int) bool {
	var idx [4]int
	var cnt [4]int
	size := 0

	add := func(v int) bool {
		for i := 0; i < size; i++ {
			if idx[i] == v {
				cnt[i]++
				if cnt[i] > avail[v] {
					return false
				}
				return true
			}
		}
		if size == len(idx) {
			return false
		}
		idx[size] = v
		cnt[size] = 1
		if cnt[size] > avail[v] {
			return false
		}
		size++
		return true
	}

	if !add(a.i) {
		return false
	}
	if !add(a.j) {
		return false
	}
	if !add(b.i) {
		return false
	}
	if !add(b.j) {
		return false
	}

	return true
}
