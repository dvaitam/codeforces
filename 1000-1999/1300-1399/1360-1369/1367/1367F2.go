package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}

	// Build unique values and per-value position lists
	sorted := append([]int{}, a...)
	sort.Ints(sorted)
	uniq := []int{sorted[0]}
	for _, v := range sorted[1:] {
		if v != uniq[len(uniq)-1] {
			uniq = append(uniq, v)
		}
	}
	m := len(uniq)
	rank := make(map[int]int, m)
	for i, v := range uniq {
		rank[v] = i
	}
	pos := make([][]int, m)
	for i, v := range a {
		pos[rank[v]] = append(pos[rank[v]], i)
	}

	maxKept := 0

	// Single-value windows
	for v := 0; v < m; v++ {
		if len(pos[v]) > maxKept {
			maxKept = len(pos[v])
		}
	}

	// Adjacent-pair windows: sweep to find optimal split
	for v := 0; v < m-1; v++ {
		pL, pR := pos[v], pos[v+1]
		curL, curR := 0, len(pR)
		if curL+curR > maxKept {
			maxKept = curL + curR
		}
		i, j := 0, 0
		for i < len(pL) || j < len(pR) {
			if i < len(pL) && (j == len(pR) || pL[i] < pR[j]) {
				curL++
				i++
			} else {
				curR--
				j++
			}
			if curL+curR > maxKept {
				maxKept = curL + curR
			}
		}
	}

	// Multi-value windows [li, ri] with ri >= li+2
	// Mandatory elements: all elements with unique index in (li, ri)
	// They must form a non-decreasing sequence when sorted by position.
	// L-optional: pos[li] elements with position < minMandatoryPos
	// R-optional: pos[ri] elements with position > maxMandatoryPos
	type posVal struct{ p, v int }
	for li := 0; li < m; li++ {
		// Build mandatory incrementally as ri increases
		mandatory := []posVal{}
		for ri := li + 2; ri < m; ri++ {
			// Add elements of value uniq[ri-1] to mandatory
			for _, p := range pos[ri-1] {
				mandatory = append(mandatory, posVal{p, uniq[ri-1]})
			}
			// Re-sort by position (new elements from pos[ri-1] are already sorted,
			// but we need to merge with existing)
			sort.Slice(mandatory, func(a, b int) bool { return mandatory[a].p < mandatory[b].p })

			// Check mandatory is non-decreasing in value
			valid := true
			for k := 1; k < len(mandatory); k++ {
				if mandatory[k].v < mandatory[k-1].v {
					valid = false
					break
				}
			}
			if !valid {
				break // adding more values won't help
			}

			minP := mandatory[0].p
			maxP := mandatory[len(mandatory)-1].p

			lOpt := sort.Search(len(pos[li]), func(k int) bool { return pos[li][k] >= minP })
			rOpt := len(pos[ri]) - sort.Search(len(pos[ri]), func(k int) bool { return pos[ri][k] > maxP })

			total := lOpt + len(mandatory) + rOpt
			if total > maxKept {
				maxKept = total
			}
		}
	}

	fmt.Fprintln(writer, n-maxKept)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}
