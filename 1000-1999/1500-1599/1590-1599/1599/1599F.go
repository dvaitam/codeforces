package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	cities := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cities[i])
	}
	// store indices for each difference between consecutive cities
	diffPos := make(map[int64][]int)
	for i := 0; i < n-1; i++ {
		d := (cities[i+1] - cities[i] + mod) % mod
		diffPos[d] = append(diffPos[d], i+1) // difference at index i+1 (1-based)
	}

	for ; q > 0; q-- {
		var l, r int
		var d int64
		fmt.Fscan(in, &l, &r, &d)
		if l >= r {
			fmt.Fprintln(out, "Yes")
			continue
		}
		positions := diffPos[d]
		if len(positions) == 0 {
			fmt.Fprintln(out, "No")
			continue
		}
		left := sort.SearchInts(positions, l)
		right := sort.SearchInts(positions, r)
		if right-left == r-l {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
