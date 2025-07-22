package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func prefixSum(a []int) []int {
	ps := make([]int, len(a)+1)
	for i, v := range a {
		ps[i+1] = ps[i] + v
	}
	return ps
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	beauty := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &beauty[i])
	}
	var colors string
	fmt.Fscan(in, &colors)

	var R, O, W []int
	for i := 0; i < n; i++ {
		switch colors[i] {
		case 'R':
			R = append(R, beauty[i])
		case 'O':
			O = append(O, beauty[i])
		case 'W':
			W = append(W, beauty[i])
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(R)))
	sort.Sort(sort.Reverse(sort.IntSlice(O)))
	sort.Sort(sort.Reverse(sort.IntSlice(W)))

	preR := prefixSum(R)
	preO := prefixSum(O)
	preW := prefixSum(W)

	best := -1

	for r := 1; r <= len(R) && r <= k-1; r++ {
		o := k - r
		if o >= 1 && o <= len(O) {
			sum := preR[r] + preO[o]
			if sum > best {
				best = sum
			}
		}
	}

	for w := 1; w <= len(W) && w <= k-1; w++ {
		o := k - w
		if o >= 1 && o <= len(O) {
			sum := preW[w] + preO[o]
			if sum > best {
				best = sum
			}
		}
	}

	fmt.Fprintln(out, best)
}
