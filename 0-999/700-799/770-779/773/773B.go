package main

import (
	"bufio"
	"fmt"
	"os"
)

func score(x int, t int) int {
	if t == -1 {
		return 0
	}
	return x * (250 - t) / 250
}
func floorMul(n, num, den int) int { return n * num / den }
func feasible(nPrime, solved, k, lnum, lden, rnum, rden int) bool {
	upper := floorMul(nPrime, rnum, rden) - solved
	if upper > k {
		upper = k
	}
	if upper < 0 {
		return false
	}
	lower := 0
	if lnum != 0 {
		limit := floorMul(nPrime, lnum, lden)
		if solved <= limit {
			lower = limit - solved + 1
		}
	}
	if lower > k {
		return false
	}
	return lower <= upper
}
func bestDiff(nPrime, solved, k int, tV, tP int) int {
	xVals := []int{500, 1000, 1500, 2000, 2500, 3000}
	lnum := []int{1, 1, 1, 1, 1, 0}
	lden := []int{2, 4, 8, 16, 32, 1}
	rnum := []int{1, 1, 1, 1, 1, 1}
	rden := []int{1, 2, 4, 8, 16, 32}
	best := -1 << 60
	if tV == -1 {
		// ratio fixed at solved/nPrime
		x := 0
		if 2*solved > nPrime {
			x = 500
		} else if 4*solved > nPrime {
			x = 1000
		} else if 8*solved > nPrime {
			x = 1500
		} else if 16*solved > nPrime {
			x = 2000
		} else if 32*solved > nPrime {
			x = 2500
		} else {
			x = 3000
		}
		diff := -score(x, tP)
		return diff
	}
	for idx := 0; idx < 6; idx++ {
		x := xVals[idx]
		ok := false
		if idx == 5 {
			if solved <= floorMul(nPrime, 1, 32) {
				ok = true
			}
		} else {
			ok = feasible(nPrime, solved, k, lnum[idx], lden[idx], rnum[idx], rden[idx])
		}
		if ok {
			diff := score(x, tV) - score(x, tP)
			if diff > best {
				best = diff
			}
		}
	}
	return best
}
func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	times := make([][]int, n)
	for i := 0; i < n; i++ {
		times[i] = make([]int, 5)
		for j := 0; j < 5; j++ {
			fmt.Fscan(in, &times[i][j])
		}
	}
	solved := make([]int, 5)
	for j := 0; j < 5; j++ {
		cnt := 0
		for i := 0; i < n; i++ {
			if times[i][j] != -1 {
				cnt++
			}
		}
		solved[j] = cnt
	}
	maxK := 5000
	for k := 0; k <= maxK; k++ {
		nPrime := n + k
		total := 0
		for j := 0; j < 5; j++ {
			diff := bestDiff(nPrime, solved[j], k, times[0][j], times[1][j])
			total += diff
		}
		if total > 0 {
			fmt.Println(k)
			return
		}
	}
	fmt.Println(-1)
}
