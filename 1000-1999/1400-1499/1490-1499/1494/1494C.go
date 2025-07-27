package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// solve computes maximum number of specials that can be matched on a single side
func solve(a, b []int) int {
	res := 0
	j := 0
	for i := 0; i < len(b); i++ {
		for j < len(a) && a[j] <= b[i] {
			j++
		}
		x := b[i] - j + 1
		idx := sort.Search(len(b), func(k int) bool { return b[k] >= x })
		if cand := i - idx + 1; cand > res {
			if cand > 0 {
				res = cand
			}
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		boxes := make([]int, n)
		specials := make([]int, m)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &boxes[i])
		}
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &specials[i])
		}
		var posA, posB, negA, negB []int
		for _, v := range boxes {
			if v > 0 {
				posA = append(posA, v)
			} else {
				negA = append(negA, -v)
			}
		}
		for _, v := range specials {
			if v > 0 {
				posB = append(posB, v)
			} else {
				negB = append(negB, -v)
			}
		}
		sort.Ints(posA)
		sort.Ints(posB)
		sort.Ints(negA)
		sort.Ints(negB)
		// reverse negatives to have increasing order from zero outward
		for i, j := 0, len(negA)-1; i < j; i, j = i+1, j-1 {
			negA[i], negA[j] = negA[j], negA[i]
		}
		for i, j := 0, len(negB)-1; i < j; i, j = i+1, j-1 {
			negB[i], negB[j] = negB[j], negB[i]
		}
		ans := solve(posA, posB) + solve(negA, negB)
		fmt.Println(ans)
	}
}
