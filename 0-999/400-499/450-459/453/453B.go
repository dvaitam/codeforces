package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	F       [60][]int
	IFm     [60][]int
	p       = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	sol     int
	am      [100]int
	v       []int
	n       int
	inArr   []int
	in2     []int
	bestDif = 30 * 100
	num     []int
	dp      [101][20]int
	bestSol []int
)

// solveDP returns minimal cost for first i elements and j candidates
func solveDP(i, j int) int {
	if i == 0 {
		return 0
	}
	if dp[i][j] != 0 {
		return dp[i][j] - 1
	}
	// skip current: keep original inArr[i-1]
	ret := solveDP(i-1, j) + inArr[i-1] - 1
	if j > 0 {
		// delete operation
		if v := solveDP(i, j-1); v < ret {
			ret = v
		}
		// replace with num[j-1]
		if v := solveDP(i-1, j-1) + abs(inArr[i-1]-num[j-1]); v < ret {
			ret = v
		}
	}
	dp[i][j] = ret + 1
	return ret
}

func solveBlock() {
	num = num[:0]
	// skip v[0] index usage; v holds chosen numbers starting at index 1
	// include all chosen values in v
	for _, val := range v {
		num = append(num, val)
	}
	// append extra primes
	num = append(num, 31, 37, 41, 43, 47, 53, 59)
	sort.Ints(num)
	// reset dp
	for i := range dp {
		for j := range dp[i] {
			dp[i][j] = 0
		}
	}
	now := solveDP(n, len(num))
	if now < bestDif {
		bestDif = now
		// reconstruct solution
		i, j := n, len(num)
		r := make([]int, 0, n)
		for i > 0 {
			if dp[i][j]-1 == solveDP(i-1, j)+inArr[i-1]-1 {
				r = append(r, 1)
				i--
			} else if j > 0 && dp[i][j]-1 == solveDP(i, j-1) {
				j--
			} else {
				r = append(r, num[j-1])
				i--
				j--
			}
		}
		// reverse r
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		bestSol = make([]int, len(r))
		copy(bestSol, r)
	}
}

func backTrack(k int) {
	if k == len(p) {
		sol++
		solveBlock()
		return
	}
	atLeastOne := false
	prime := p[k]
	for _, now := range IFm[prime] {
		ok := true
		for _, f := range F[now] {
			if am[f] != 0 {
				ok = false
				break
			}
		}
		if ok {
			atLeastOne = true
			for _, f := range F[now] {
				am[f] = 1
			}
			v = append(v, now)
			backTrack(k + 1)
			v = v[:len(v)-1]
			for _, f := range F[now] {
				am[f] = 0
			}
		}
	}
	if !atLeastOne {
		backTrack(k + 1)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	inArr = make([]int, n)
	in2 = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &inArr[i])
		in2[i] = inArr[i]
	}
	sort.Ints(inArr)
	// build factors
	for i := 2; i < 60; i++ {
		x := i
		for f := 2; f <= x; f++ {
			if x%f == 0 {
				F[i] = append(F[i], f)
				IFm[f] = append(IFm[f], i)
				for x%f == 0 {
					x /= f
				}
			}
		}
	}
	// init v with length 1 dummy for 1-based index
	v = make([]int, 0, 10)
	backTrack(0)
	// output bestSol in original order
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		// find original position in sorted inArr
		for j := 0; j < n; j++ {
			if in2[i] == inArr[j] {
				inArr[j] = -1
				fmt.Fprint(out, bestSol[j])
				break
			}
		}
	}
	fmt.Fprintln(out)
}
