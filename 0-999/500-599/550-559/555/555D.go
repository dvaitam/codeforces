package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	x   int64
	idx int
}

var xs []int64
var pos []int
var n int

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func getNext(dir int, idx int, length int64) int {
	if dir == 1 { // move to the right
		target := xs[idx] + length
		j := sort.Search(len(xs), func(i int) bool { return xs[i] > target }) - 1
		if j < idx {
			j = idx
		}
		return j
	}
	// dir == 0, move left
	target := xs[idx] - length
	j := sort.Search(len(xs), func(i int) bool { return xs[i] >= target })
	if j > idx {
		j = idx
	}
	return j
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	var m int
	fmt.Fscan(in, &m)
	pegs := make([]pair, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pegs[i].x)
		pegs[i].idx = i
	}
	sort.Slice(pegs, func(i, j int) bool { return pegs[i].x < pegs[j].x })
	xs = make([]int64, n)
	pos = make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = pegs[i].x
		pos[pegs[i].idx] = i
	}
	out := bufio.NewWriter(os.Stdout)
	for ; m > 0; m-- {
		var a int
		var l int64
		fmt.Fscan(in, &a, &l)
		a--
		idx := pos[a]
		dir := 1
		prev := [2]int{-1, -1}
		for {
			nextIdx := getNext(dir, idx, l)
			if nextIdx == idx && nextIdx == prev[dir] {
				break
			}
			if nextIdx == prev[dir] {
				dist := abs(xs[prev[dir]]-xs[prev[dir^1]]) * 2
				if dist > 0 {
					l %= dist
				}
				nextIdx = getNext(dir, idx, l)
			}
			prev[dir] = nextIdx
			l -= abs(xs[nextIdx] - xs[idx])
			idx = nextIdx
			dir ^= 1
		}
		fmt.Fprintln(out, pegs[idx].idx+1)
	}
	out.Flush()
}
