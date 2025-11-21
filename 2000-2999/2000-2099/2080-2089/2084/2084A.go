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

	var t int
	fmt.Fscan(in, &t)
	maxN := 100
	precomputed := make(map[int][]int)
	precomputed[3] = []int{3, 1, 2}
	precomputed[5] = []int{5, 2, 1, 3, 4}
	precomputed[7] = []int{7, 2, 1, 3, 4, 5, 6}
	precomputed[9] = []int{9, 5, 1, 2, 3, 4, 6, 7, 8}
	for i := 1; i <= maxN; i += 2 {
		if _, ok := precomputed[i]; !ok && i >= 3 {
			precomputed[i] = buildPermutation(i)
		}
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n%2 == 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		perm := precomputed[n]
		if perm == nil {
			perm = buildPermutation(n)
			precomputed[n] = perm
		}
		for i, v := range perm {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

func buildPermutation(n int) []int {
	perm := make([]int, n)
	perm[0] = n
	idx := 1
	for i := 2; idx < n; i += 2 {
		perm[idx] = i
		idx++
	}
	for i := 1; idx < n; i += 2 {
		perm[idx] = i
		idx++
	}
	return perm
}
