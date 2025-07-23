package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for CF 808E - Selling Souvenirs
// using greedy merging of weight2 items and pairs of weight1 items.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var w1, w2, w3 []int64
	for i := 0; i < n; i++ {
		var w int
		var c int64
		fmt.Fscan(reader, &w, &c)
		switch w {
		case 1:
			w1 = append(w1, c)
		case 2:
			w2 = append(w2, c)
		case 3:
			w3 = append(w3, c)
		}
	}

	sort.Slice(w1, func(i, j int) bool { return w1[i] > w1[j] })
	sort.Slice(w2, func(i, j int) bool { return w2[i] > w2[j] })
	sort.Slice(w3, func(i, j int) bool { return w3[i] > w3[j] })

	pre3 := make([]int64, len(w3)+1)
	for i, v := range w3 {
		pre3[i+1] = pre3[i] + v
	}

	// form pairs from weight1 items
	pairs := make([]int64, 0, len(w1)/2)
	for i := 0; i+1 < len(w1); i += 2 {
		pairs = append(pairs, w1[i]+w1[i+1])
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i] > pairs[j] })

	// merge weight2 items and pairs into a single sorted list of weight=2 items
	maxTwoItems := len(w2) + len(pairs)
	if maxTwoItems > m/2 {
		maxTwoItems = m / 2
	}
	twoCosts := make([]int64, 0, maxTwoItems)
	pairFlag := make([]int, 0, maxTwoItems)
	i2, ip := 0, 0
	for len(twoCosts) < maxTwoItems {
		if i2 < len(w2) && (ip >= len(pairs) || w2[i2] >= pairs[ip]) {
			twoCosts = append(twoCosts, w2[i2])
			pairFlag = append(pairFlag, 0)
			i2++
		} else if ip < len(pairs) {
			twoCosts = append(twoCosts, pairs[ip])
			pairFlag = append(pairFlag, 1)
			ip++
		} else {
			break
		}
	}

	preTwo := make([]int64, len(twoCosts)+1)
	prePairs := make([]int, len(twoCosts)+1)
	for i, v := range twoCosts {
		preTwo[i+1] = preTwo[i] + v
		prePairs[i+1] = prePairs[i] + pairFlag[i]
	}

	best := int64(0)
	maxT3 := m / 3
	if maxT3 > len(w3) {
		maxT3 = len(w3)
	}
	for t3 := 0; t3 <= maxT3; t3++ {
		remaining := m - 3*t3
		cost := pre3[t3]

		t2 := remaining / 2
		if t2 > len(twoCosts) {
			t2 = len(twoCosts)
		}
		cost += preTwo[t2]
		usedPairs := prePairs[t2]
		idx := usedPairs * 2
		if remaining%2 == 1 && idx < len(w1) {
			cost += w1[idx]
		}
		if cost > best {
			best = cost
		}
	}

	fmt.Fprintln(writer, best)
}
